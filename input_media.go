package gotdbot

type InputMediaOpts struct {
	Caption               string
	CaptionEntities       []TextEntity
	ParseMode             string
	HasSpoiler            bool
	ShowCaptionAboveMedia bool
	Thumbnail             *InputThumbnail
	Width                 int32
	Height                int32
	Duration              int32
	SupportsStreaming     bool
	Performer             string
	Title                 string
	DisableTypeDetection  bool
	AddedStickerFileIds   []int32
	Cover                 InputFile
	StartTimestamp        int32
}

func (c *Client) captionText(opts *InputMediaOpts) (*FormattedText, error) {
	if opts == nil {
		return &FormattedText{}, nil
	}
	return c.GetFormattedText(opts.Caption, opts.CaptionEntities, opts.ParseMode)
}

func (c *Client) InputMediaPhoto(file InputFile, opts *InputMediaOpts) (InputMessageContent, error) {
	if opts == nil {
		opts = &InputMediaOpts{}
	}
	caption, err := c.captionText(opts)
	if err != nil {
		return nil, err
	}
	return &InputMessagePhoto{
		Photo: &InputPhoto{
			Photo:               file,
			Thumbnail:           opts.Thumbnail,
			AddedStickerFileIds: opts.AddedStickerFileIds,
			Width:               opts.Width,
			Height:              opts.Height,
		},
		Caption:               caption,
		HasSpoiler:            opts.HasSpoiler,
		ShowCaptionAboveMedia: opts.ShowCaptionAboveMedia,
	}, nil
}

func (c *Client) InputMediaVideo(file InputFile, opts *InputMediaOpts) (InputMessageContent, error) {
	if opts == nil {
		opts = &InputMediaOpts{}
	}
	caption, err := c.captionText(opts)
	if err != nil {
		return nil, err
	}
	return &InputMessageVideo{
		Video: &InputVideo{
			Video:               file,
			Thumbnail:           opts.Thumbnail,
			AddedStickerFileIds: opts.AddedStickerFileIds,
			Duration:            opts.Duration,
			Width:               opts.Width,
			Height:              opts.Height,
			SupportsStreaming:   opts.SupportsStreaming,
			StartTimestamp:      opts.StartTimestamp,
			Cover:               opts.Cover,
		},
		Caption:               caption,
		HasSpoiler:            opts.HasSpoiler,
		ShowCaptionAboveMedia: opts.ShowCaptionAboveMedia,
	}, nil
}

func (c *Client) InputMediaAnimation(file InputFile, opts *InputMediaOpts) (InputMessageContent, error) {
	if opts == nil {
		opts = &InputMediaOpts{}
	}
	caption, err := c.captionText(opts)
	if err != nil {
		return nil, err
	}
	return &InputMessageAnimation{
		Animation: &InputAnimation{
			Animation:           file,
			Thumbnail:           opts.Thumbnail,
			AddedStickerFileIds: opts.AddedStickerFileIds,
			Duration:            opts.Duration,
			Width:               opts.Width,
			Height:              opts.Height,
		},
		Caption:               caption,
		HasSpoiler:            opts.HasSpoiler,
		ShowCaptionAboveMedia: opts.ShowCaptionAboveMedia,
	}, nil
}

func (c *Client) InputMediaAudio(file InputFile, opts *InputMediaOpts) (InputMessageContent, error) {
	if opts == nil {
		opts = &InputMediaOpts{}
	}
	caption, err := c.captionText(opts)
	if err != nil {
		return nil, err
	}
	return &InputMessageAudio{
		Audio: &InputAudio{
			Audio:               file,
			Duration:            opts.Duration,
			Title:               opts.Title,
			Performer:           opts.Performer,
			AlbumCoverThumbnail: opts.Thumbnail,
		},
		Caption: caption,
	}, nil
}

func (c *Client) InputMediaDocument(file InputFile, opts *InputMediaOpts) (InputMessageContent, error) {
	if opts == nil {
		opts = &InputMediaOpts{}
	}
	caption, err := c.captionText(opts)
	if err != nil {
		return nil, err
	}
	return &InputMessageDocument{
		Document: &InputDocument{
			Document:                    file,
			Thumbnail:                   opts.Thumbnail,
			DisableContentTypeDetection: opts.DisableTypeDetection,
		},
		Caption: caption,
	}, nil
}

func InputMediaPhotoPath(path string) InputMessageContent {
	return &InputMessagePhoto{
		Photo: &InputPhoto{Photo: FileFromPath(path)},
	}
}

func InputMediaVideoPath(path string) InputMessageContent {
	return &InputMessageVideo{
		Video: &InputVideo{Video: FileFromPath(path), SupportsStreaming: true},
	}
}

func InputMediaDocumentPath(path string) InputMessageContent {
	return &InputMessageDocument{
		Document: &InputDocument{Document: FileFromPath(path)},
	}
}

func InputMediaAudioPath(path string) InputMessageContent {
	return &InputMessageAudio{
		Audio: &InputAudio{Audio: FileFromPath(path)},
	}
}

func InputMediaAnimationPath(path string) InputMessageContent {
	return &InputMessageAnimation{
		Animation: &InputAnimation{Animation: FileFromPath(path)},
	}
}
