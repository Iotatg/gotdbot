package gotdbot

import (
	"context"
	"time"
)

type DownloadProgressFunc func(file *File)

type DownloadProgressOpts struct {
	Priority   int32
	Offset     int64
	Limit      int64
	Timeout    time.Duration
	OnProgress DownloadProgressFunc
}

func (c *Client) DownloadFileProgress(fileId int32, opts *DownloadProgressOpts) (*File, error) {
	if opts == nil {
		opts = &DownloadProgressOpts{}
	}
	priority := opts.Priority
	if priority <= 0 {
		priority = 1
	}
	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Minute
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	progressCh := make(chan *File, 8)
	h := &updateFileHandler{
		filter: func(u *UpdateFile) bool {
			return u != nil && u.File != nil && u.File.Id == fileId
		},
		response: func(client *Client, update *UpdateFile) error {
			select {
			case progressCh <- update.File:
			default:
			}
			return ContinueHandlers
		},
	}
	c.AddHandlerGroup(h, -90)
	defer c.RemoveHandlerGroup(h, -90)

	started, err := c.DownloadFile(fileId, opts.Limit, opts.Offset, priority, &DownloadFileOpts{Synchronous: false})
	if err != nil {
		return nil, err
	}
	if started != nil && started.Local != nil && started.Local.IsDownloadingCompleted {
		if opts.OnProgress != nil {
			opts.OnProgress(started)
		}
		return started, nil
	}

	var latest *File
	if started != nil {
		latest = started
		if opts.OnProgress != nil {
			opts.OnProgress(started)
		}
	}

	for {
		select {
		case <-ctx.Done():
			if latest != nil && latest.Local != nil && latest.Local.IsDownloadingCompleted {
				return latest, nil
			}
			if ctx.Err() == context.DeadlineExceeded {
				return latest, ErrDownloadStopped
			}
			return latest, ctx.Err()
		case f := <-progressCh:
			latest = f
			if opts.OnProgress != nil {
				opts.OnProgress(f)
			}
			if f != nil && f.Local != nil && f.Local.IsDownloadingCompleted {
				return f, nil
			}
		}
	}
}

func (m *Message) DownloadProgress(c *Client, opts *DownloadProgressOpts) (*File, error) {
	f, err := m.Download(c, 1, 0, 0, false)
	if err != nil || f == nil {
		return f, err
	}
	if f.Local != nil && f.Local.IsDownloadingCompleted {
		if opts != nil && opts.OnProgress != nil {
			opts.OnProgress(f)
		}
		return f, nil
	}
	return c.DownloadFileProgress(f.Id, opts)
}

func (f *File) ProgressPercent() float64 {
	if f == nil {
		return 0
	}
	total := f.Size
	if total <= 0 {
		total = f.ExpectedSize
	}
	if total <= 0 || f.Local == nil {
		return 0
	}
	return float64(f.Local.DownloadedSize) * 100 / float64(total)
}
