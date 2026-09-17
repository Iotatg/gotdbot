package gotdbot

func (m *Message) Photo() *Photo {
	if m == nil {
		return nil
	}
	c, ok := m.Content.(*MessagePhoto)
	if !ok {
		return nil
	}
	return c.Photo
}

func (m *Message) Video() *Video {
	if m == nil {
		return nil
	}
	c, ok := m.Content.(*MessageVideo)
	if !ok {
		return nil
	}
	return c.Video
}

func (m *Message) Animation() *Animation {
	if m == nil {
		return nil
	}
	c, ok := m.Content.(*MessageAnimation)
	if !ok {
		return nil
	}
	return c.Animation
}

func (m *Message) Audio() *Audio {
	if m == nil {
		return nil
	}
	c, ok := m.Content.(*MessageAudio)
	if !ok {
		return nil
	}
	return c.Audio
}

func (m *Message) Document() *Document {
	if m == nil {
		return nil
	}
	c, ok := m.Content.(*MessageDocument)
	if !ok {
		return nil
	}
	return c.Document
}

func (m *Message) Sticker() *Sticker {
	if m == nil {
		return nil
	}
	c, ok := m.Content.(*MessageSticker)
	if !ok {
		return nil
	}
	return c.Sticker
}

func (m *Message) Voice() *VoiceNote {
	if m == nil {
		return nil
	}
	c, ok := m.Content.(*MessageVoiceNote)
	if !ok {
		return nil
	}
	return c.VoiceNote
}

func (m *Message) VideoNote() *VideoNote {
	if m == nil {
		return nil
	}
	c, ok := m.Content.(*MessageVideoNote)
	if !ok {
		return nil
	}
	return c.VideoNote
}

func (m *Message) Contact() *Contact {
	if m == nil {
		return nil
	}
	c, ok := m.Content.(*MessageContact)
	if !ok {
		return nil
	}
	return c.Contact
}

func (m *Message) Location() *Location {
	if m == nil {
		return nil
	}
	c, ok := m.Content.(*MessageLocation)
	if !ok {
		return nil
	}
	return c.Location
}

func (m *Message) Venue() *Venue {
	if m == nil {
		return nil
	}
	c, ok := m.Content.(*MessageVenue)
	if !ok {
		return nil
	}
	return c.Venue
}

func (m *Message) Poll() *Poll {
	if m == nil {
		return nil
	}
	c, ok := m.Content.(*MessagePoll)
	if !ok {
		return nil
	}
	return c.Poll
}

func (m *Message) Dice() *MessageDice {
	if m == nil {
		return nil
	}
	c, ok := m.Content.(*MessageDice)
	if !ok {
		return nil
	}
	return c
}

func (m *Message) Game() *Game {
	if m == nil {
		return nil
	}
	c, ok := m.Content.(*MessageGame)
	if !ok {
		return nil
	}
	return c.Game
}

func (m *Message) Story() *MessageStory {
	if m == nil {
		return nil
	}
	c, ok := m.Content.(*MessageStory)
	if !ok {
		return nil
	}
	return c
}

func (m *Message) Invoice() *MessageInvoice {
	if m == nil {
		return nil
	}
	c, ok := m.Content.(*MessageInvoice)
	if !ok {
		return nil
	}
	return c
}

func (m *Message) PaidMedia() *MessagePaidMedia {
	if m == nil {
		return nil
	}
	c, ok := m.Content.(*MessagePaidMedia)
	if !ok {
		return nil
	}
	return c
}

func (m *Message) Checklist() *MessageChecklist {
	if m == nil {
		return nil
	}
	c, ok := m.Content.(*MessageChecklist)
	if !ok {
		return nil
	}
	return c
}

func (m *Message) Gift() *MessageGift {
	if m == nil {
		return nil
	}
	c, ok := m.Content.(*MessageGift)
	if !ok {
		return nil
	}
	return c
}

func (m *Message) Giveaway() *MessageGiveaway {
	if m == nil {
		return nil
	}
	c, ok := m.Content.(*MessageGiveaway)
	if !ok {
		return nil
	}
	return c
}

func (m *Message) GiveawayWinners() *MessageGiveawayWinners {
	if m == nil {
		return nil
	}
	c, ok := m.Content.(*MessageGiveawayWinners)
	if !ok {
		return nil
	}
	return c
}

func (m *Message) UsersShared() *MessageUsersShared {
	if m == nil {
		return nil
	}
	c, ok := m.Content.(*MessageUsersShared)
	if !ok {
		return nil
	}
	return c
}

func (m *Message) ChatShared() *MessageChatShared {
	if m == nil {
		return nil
	}
	c, ok := m.Content.(*MessageChatShared)
	if !ok {
		return nil
	}
	return c
}

func (m *Message) WebAppData() *MessageWebAppDataReceived {
	if m == nil {
		return nil
	}
	c, ok := m.Content.(*MessageWebAppDataReceived)
	if !ok {
		return nil
	}
	return c
}

func (m *Message) HasPhoto() bool       { return m.Photo() != nil }
func (m *Message) HasVideo() bool       { return m.Video() != nil }
func (m *Message) HasAnimation() bool   { return m.Animation() != nil }
func (m *Message) HasAudio() bool       { return m.Audio() != nil }
func (m *Message) HasDocument() bool    { return m.Document() != nil }
func (m *Message) HasSticker() bool     { return m.Sticker() != nil }
func (m *Message) HasVoice() bool       { return m.Voice() != nil }
func (m *Message) HasVideoNote() bool   { return m.VideoNote() != nil }
func (m *Message) HasContact() bool     { return m.Contact() != nil }
func (m *Message) HasLocation() bool    { return m.Location() != nil }
func (m *Message) HasVenue() bool       { return m.Venue() != nil }
func (m *Message) HasPoll() bool        { return m.Poll() != nil }
func (m *Message) HasDice() bool        { return m.Dice() != nil }
func (m *Message) HasGame() bool        { return m.Game() != nil }
func (m *Message) HasStory() bool       { return m.Story() != nil }
func (m *Message) HasInvoice() bool     { return m.Invoice() != nil }
func (m *Message) HasPaidMedia() bool   { return m.PaidMedia() != nil }
func (m *Message) HasChecklist() bool   { return m.Checklist() != nil }
func (m *Message) HasGift() bool        { return m.Gift() != nil }
func (m *Message) HasGiveaway() bool    { return m.Giveaway() != nil }
func (m *Message) HasText() bool        { return m.Text() != "" }
func (m *Message) HasCaption() bool     { return m.Caption() != "" }
func (m *Message) HasReplyMarkup() bool { return m != nil && m.ReplyMarkup != nil }

func (m *Message) HasMedia() bool {
	if m == nil {
		return false
	}
	return m.HasPhoto() || m.HasVideo() || m.HasAnimation() || m.HasAudio() ||
		m.HasDocument() || m.HasSticker() || m.HasVoice() || m.HasVideoNote() ||
		m.HasPaidMedia()
}

func (m *Message) HasSpoiler() bool {
	if m == nil || m.Content == nil {
		return false
	}
	switch c := m.Content.(type) {
	case *MessagePhoto:
		return c.HasSpoiler
	case *MessageVideo:
		return c.HasSpoiler
	case *MessageAnimation:
		return c.HasSpoiler
	default:
		return false
	}
}

func (m *Message) IsForwarded() bool {
	return m != nil && m.ForwardInfo != nil
}

func (m *Message) IsReply() bool {
	return m != nil && m.ReplyTo != nil
}

func (m *Message) IsAlbum() bool {
	return m != nil && m.MediaAlbumId != 0
}

func (m *Message) IsViaBot() bool {
	return m != nil && m.ViaBotUserId != 0
}

func (m *Message) IsScheduled() bool {
	return m != nil && m.SchedulingState != nil
}

func (m *Message) IsChannelPostMessage() bool {
	return m != nil && m.IsChannelPost
}

func (m *Message) IsForum() bool {
	if m == nil || m.TopicId == nil {
		return false
	}
	_, ok := m.TopicId.(*MessageTopicForum)
	return ok
}

func (m *Message) IsDirect() bool {
	if m == nil || m.TopicId == nil {
		return false
	}
	_, ok := m.TopicId.(*MessageTopicDirectMessages)
	return ok
}

func (m *Message) ForumTopicID() int32 {
	if m == nil || m.TopicId == nil {
		return 0
	}
	if t, ok := m.TopicId.(*MessageTopicForum); ok {
		return t.ForumTopicId
	}
	return 0
}

func (m *Message) IsService() bool {
	if m == nil || m.Content == nil {
		return false
	}
	switch m.Content.(type) {
	case *MessageText, *MessagePhoto, *MessageVideo, *MessageAnimation, *MessageAudio,
		*MessageDocument, *MessageSticker, *MessageVoiceNote, *MessageVideoNote,
		*MessageContact, *MessageLocation, *MessageVenue, *MessagePoll, *MessageDice,
		*MessageGame, *MessageStory, *MessageInvoice, *MessagePaidMedia, *MessageChecklist,
		*MessageGift, *MessageGiveaway, *MessageGiveawayWinners:
		return false
	default:
		return true
	}
}
