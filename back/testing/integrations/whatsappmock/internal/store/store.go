package store

import (
	"sync"
	"time"
)

type Send struct {
	MessageID     string            `json:"message_id"`
	PhoneNumberID string            `json:"phone_number_id"`
	To            string            `json:"to"`
	TemplateName  string            `json:"template_name"`
	Language      string            `json:"language"`
	Parameters    []string          `json:"parameters"`
	Buttons       []string          `json:"buttons"`
	Variables     map[string]string `json:"variables"`
	SentAt        time.Time         `json:"sent_at"`
	RepliedWith   string            `json:"replied_with"`
}

type Template struct {
	MetaID   string    `json:"meta_id"`
	WABAID   string    `json:"waba_id"`
	Name     string    `json:"name"`
	Language string    `json:"language"`
	Status   string    `json:"status"`
	SentAt   time.Time `json:"created_at"`
}

type Store struct {
	mu        sync.RWMutex
	sends     []*Send
	templates []*Template
	script    map[string]string
	autoReply bool
}

func New() *Store {
	return &Store{
		script:    defaultScript(),
		autoReply: true,
	}
}

func defaultScript() map[string]string {
	return map[string]string{
		"26_saludo_inicial": "Si",
		"26_dice_que_no":    "eres una rata",
		"26_hola":           "Jodete",
		"26_me_jodo":        "genial",
	}
}

func (s *Store) AddSend(send *Send) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sends = append(s.sends, send)
}

func (s *Store) MarkReplied(messageID, buttonText string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, send := range s.sends {
		if send.MessageID == messageID {
			send.RepliedWith = buttonText
			return
		}
	}
}

func (s *Store) Sends() []*Send {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Send, len(s.sends))
	copy(out, s.sends)
	return out
}

func (s *Store) LastSendTo(phone string) *Send {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := len(s.sends) - 1; i >= 0; i-- {
		if s.sends[i].To == phone {
			return s.sends[i]
		}
	}
	return nil
}

func (s *Store) AddTemplate(template *Template) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.templates = append(s.templates, template)
}

func (s *Store) Templates() []*Template {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Template, len(s.templates))
	copy(out, s.templates)
	return out
}

func (s *Store) TemplateByName(name string) *Template {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := len(s.templates) - 1; i >= 0; i-- {
		if s.templates[i].Name == name {
			return s.templates[i]
		}
	}
	return nil
}

func (s *Store) ReplyFor(templateName string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !s.autoReply {
		return "", false
	}
	reply, ok := s.script[templateName]
	if !ok || reply == "" {
		return "", false
	}
	return reply, true
}

func (s *Store) SetScript(script map[string]string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for name, reply := range script {
		s.script[name] = reply
	}
}

func (s *Store) Script() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]string, len(s.script))
	for name, reply := range s.script {
		out[name] = reply
	}
	return out
}

func (s *Store) SetAutoReply(enabled bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.autoReply = enabled
}

func (s *Store) AutoReply() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.autoReply
}

func (s *Store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sends = nil
	s.templates = nil
	s.script = defaultScript()
	s.autoReply = true
}
