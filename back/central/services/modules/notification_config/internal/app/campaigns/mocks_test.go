package campaigns

import (
	"context"
	"errors"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/log"
)

var errBoom = errors.New("boom")

type campaignRepoMock struct {
	created      *entities.Campaign
	updated      []entities.Campaign
	byID         map[uint]*entities.Campaign
	list         []entities.Campaign
	listTotal    int64
	due          []entities.Campaign
	deletedID    uint
	batchMarked  uint
	countersRun  uint
	failCreate   bool
	failUpdate   bool
	failGet      bool
	failList     bool
	failDue      bool
	failDelete   bool
	failCounters bool
	failBatch    bool
}

func newCampaignRepoMock() *campaignRepoMock {
	return &campaignRepoMock{byID: map[uint]*entities.Campaign{}}
}

func (m *campaignRepoMock) CreateCampaign(_ context.Context, campaign *entities.Campaign) error {
	if m.failCreate {
		return errBoom
	}
	campaign.ID = 1
	m.created = campaign
	m.byID[campaign.ID] = campaign
	return nil
}

func (m *campaignRepoMock) UpdateCampaign(_ context.Context, campaign *entities.Campaign) error {
	if m.failUpdate {
		return errBoom
	}
	m.updated = append(m.updated, *campaign)
	m.byID[campaign.ID] = campaign
	return nil
}

func (m *campaignRepoMock) GetCampaignByID(_ context.Context, id uint) (*entities.Campaign, error) {
	if m.failGet {
		return nil, errBoom
	}
	return m.byID[id], nil
}

func (m *campaignRepoMock) ListCampaigns(_ context.Context, _ uint, _ string, _, _ int) ([]entities.Campaign, int64, error) {
	if m.failList {
		return nil, 0, errBoom
	}
	return m.list, m.listTotal, nil
}

func (m *campaignRepoMock) DeleteCampaign(_ context.Context, id uint) error {
	if m.failDelete {
		return errBoom
	}
	m.deletedID = id
	return nil
}

func (m *campaignRepoMock) ListDueCampaigns(_ context.Context, _ time.Time, _ int) ([]entities.Campaign, error) {
	if m.failDue {
		return nil, errBoom
	}
	return m.due, nil
}

func (m *campaignRepoMock) UpdateCampaignCounters(_ context.Context, campaignID uint) error {
	if m.failCounters {
		return errBoom
	}
	m.countersRun = campaignID
	return nil
}

func (m *campaignRepoMock) MarkCampaignBatch(_ context.Context, campaignID uint, _ time.Time) error {
	if m.failBatch {
		return errBoom
	}
	m.batchMarked = campaignID
	return nil
}

type sendRepoMock struct {
	bulkCreated  []entities.CampaignSend
	bulkCount    int
	pending      []entities.CampaignSend
	list         []entities.CampaignSend
	listTotal    int64
	sentSince    int64
	queued       []uint
	results      []string
	failBulk     bool
	failPending  bool
	failList     bool
	failCount    bool
	failQueued   bool
	failResult   bool
	queuedErrFor uint
}

func (m *sendRepoMock) BulkCreateSends(_ context.Context, sends []entities.CampaignSend) (int, error) {
	if m.failBulk {
		return 0, errBoom
	}
	m.bulkCreated = sends
	if m.bulkCount > 0 {
		return m.bulkCount, nil
	}
	return len(sends), nil
}

func (m *sendRepoMock) ListPendingSends(_ context.Context, _ uint, limit int) ([]entities.CampaignSend, error) {
	if m.failPending {
		return nil, errBoom
	}
	if limit < len(m.pending) {
		return m.pending[:limit], nil
	}
	return m.pending, nil
}

func (m *sendRepoMock) ListSends(_ context.Context, _ uint, _ string, _, _ int) ([]entities.CampaignSend, int64, error) {
	if m.failList {
		return nil, 0, errBoom
	}
	return m.list, m.listTotal, nil
}

func (m *sendRepoMock) MarkSendQueued(_ context.Context, sendID uint, _ time.Time) error {
	if m.failQueued || m.queuedErrFor == sendID {
		return errBoom
	}
	m.queued = append(m.queued, sendID)
	return nil
}

func (m *sendRepoMock) MarkSendResult(_ context.Context, _ uint, status, _, _ string) error {
	if m.failResult {
		return errBoom
	}
	m.results = append(m.results, status)
	return nil
}

func (m *sendRepoMock) CountSentSince(_ context.Context, _ uint, _ time.Time) (int64, error) {
	if m.failCount {
		return 0, errBoom
	}
	return m.sentSince, nil
}

type audienceMock struct {
	candidates []entities.CampaignCandidate
	total      uint
	optedOut   uint
	noPhone    uint
	reachable  uint
	failFind   bool
	failCount  bool
}

func (m *audienceMock) FindCampaignCandidates(_ context.Context, _ uint, _ entities.CampaignAudienceParams, _ string, limit int) ([]entities.CampaignCandidate, error) {
	if m.failFind {
		return nil, errBoom
	}
	if limit < len(m.candidates) {
		return m.candidates[:limit], nil
	}
	return m.candidates, nil
}

func (m *audienceMock) CountCampaignAudience(_ context.Context, _ uint, _ entities.CampaignAudienceParams, _ string) (uint, uint, uint, uint, error) {
	if m.failCount {
		return 0, 0, 0, 0, errBoom
	}
	return m.total, m.optedOut, m.noPhone, m.reachable, nil
}

type senderMock struct {
	integrationID uint
	phone         string
	fail          bool
}

func (m *senderMock) GetOwnWhatsappSender(_ context.Context, _ uint) (uint, string, error) {
	if m.fail {
		return 0, "", errBoom
	}
	return m.integrationID, m.phone, nil
}

type templateRepoMock struct {
	template *entities.WhatsappTemplate
	fail     bool
}

func (m *templateRepoMock) CreateTemplate(_ context.Context, _ *entities.WhatsappTemplate) error {
	return nil
}

func (m *templateRepoMock) UpdateTemplate(_ context.Context, _ *entities.WhatsappTemplate) error {
	return nil
}

func (m *templateRepoMock) GetTemplateByID(_ context.Context, _ uint) (*entities.WhatsappTemplate, error) {
	if m.fail {
		return nil, errBoom
	}
	return m.template, nil
}

func (m *templateRepoMock) GetTemplateByName(_ context.Context, _ uint, _, _ string) (*entities.WhatsappTemplate, error) {
	return nil, nil
}

func (m *templateRepoMock) ListTemplates(_ context.Context, _ uint, _, _ string, _, _ int) ([]entities.WhatsappTemplate, int64, error) {
	return nil, 0, nil
}

func (m *templateRepoMock) DeleteTemplate(_ context.Context, _ uint) error {
	return nil
}

func (m *templateRepoMock) UpdateTemplateStatusByMeta(_ context.Context, _, _, _, _, _ string) error {
	return nil
}

type publisherMock struct {
	published []dtos.CampaignSendMessage
	fail      bool
}

func (m *publisherMock) PublishCampaignSend(_ context.Context, message dtos.CampaignSendMessage) error {
	if m.fail {
		return errBoom
	}
	m.published = append(m.published, message)
	return nil
}

type harness struct {
	uc        IUseCase
	impl      *useCase
	campaigns *campaignRepoMock
	sends     *sendRepoMock
	audience  *audienceMock
	sender    *senderMock
	templates *templateRepoMock
	publisher *publisherMock
}

func approvedTemplate(businessID uint) *entities.WhatsappTemplate {
	return &entities.WhatsappTemplate{
		ID:         7,
		BusinessID: &businessID,
		Origin:     entities.TemplateOriginBusiness,
		Scope:      entities.TemplateScopeCampaign,
		Name:       "ruta_30_dian",
		Language:   "es",
		Status:     entities.TemplateStatusApproved,
		Variables: []entities.TemplateVariable{
			{Position: 1, Source: "customer.first_name"},
			{Position: 2, Source: "sender.name"},
		},
	}
}

func newHarness() *harness {
	businessID := uint(26)

	h := &harness{
		campaigns: newCampaignRepoMock(),
		sends:     &sendRepoMock{},
		audience:  &audienceMock{},
		sender:    &senderMock{integrationID: 99, phone: "573001112233"},
		templates: &templateRepoMock{template: approvedTemplate(businessID)},
		publisher: &publisherMock{},
	}

	h.uc = New(h.campaigns, h.sends, h.audience, h.sender, h.templates, h.publisher, log.New())
	h.impl = h.uc.(*useCase)

	return h
}
