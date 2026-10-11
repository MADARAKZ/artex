package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/artex/notify"
)

// English text（English text settings English text，English text）。
const (
	// settingNotifyEnabled English text。English text：English text，
	// English text——English text「English text」。
	settingNotifyEnabled = "notify_enabled"
	// settingNotifyPublicBaseURL English text
	// （English text https://artex.example.com）。English text。
	// English text，English text。
	settingNotifyPublicBaseURL = "notify_public_base_url"
	// settingNotifyDigestMinutes English text（English text）。
	settingNotifyDigestMinutes = "notify_digest_interval_min"
)

const (
	// notifyTick English text。3 English text，
	// English text「English text」English text「English text IM」English text。
	notifyTick = 3 * time.Second
	// notifyLease English text。English text
	// （notify English text HTTP English text 15 English text），English text
	// dispatcher English text。
	notifyLease = 3 * time.Minute
	// notifyFanOutPerTick English text，English text
	// English text。
	notifyFanOutPerTick = 200
	// notifyDefaultDigestMinutes English text。
	notifyDefaultDigestMinutes = 30
	// notifyUnlimitedBurstPerTick English text。
	// English text「English text + English text」English text
	// English text。
	notifyUnlimitedBurstPerTick = 50
	// notifyMaxSendsPerChannelPerTick English text。
	//
	// English text**English text**English text：English text（notifyLease = 3 English text），
	// English text，English text。
	// English text（Run English text goroutine English text，tick English text），English text**English text
	// English text**English text，English text，English text
	// attempts English text、English text。
	//
	// English text：3 English text / 30 English text = 6 English text**English text**、English text，
	// English text；English text 5 English text 150 English text 30 English text。English text
	// TestNotifyTickBudgetFitsWithinLease English text——English text notifyLease、
	// notifySendTimeout English text。
	notifyMaxSendsPerChannelPerTick = 5
	// notifySendTimeout English text。English text，
	// English text notifyLease，English text TestNotifyTickBudgetFitsWithinLease。
	notifySendTimeout = 30 * time.Second
)

// notifyBackoff English text，English text。
// 3 English text（English text）English text db.MaxNotifyAttempts English text，English text。
var notifyBackoff = []time.Duration{
	time.Second,
	5 * time.Second,
	30 * time.Second,
}

// Notifier English text。
//
// English text Scheduler English text，English text goroutine English text（English text server.New）。English text
// Scheduler English text tick：English text（3 English text）English text，
// English text——English text agent English text。
type Notifier struct {
	s  *Server
	pg *db.DB

	// mu English text buckets。English text、English text，English text，
	// English text。
	mu      sync.Mutex
	buckets map[int64]*notifyBucket
}

// notifyBucket English text。
//
// English text「English text」English text，English text：
// English text 20 English text、English text 20 English text，English text 40 English text，
// English text；English text，English text。
type notifyBucket struct {
	tokens   float64
	lastFill time.Time
}

func newNotifier(s *Server) *Notifier {
	return &Notifier{s: s, pg: s.m.pg, buckets: map[int64]*notifyBucket{}}
}

// Run English text ctx English text。English text server.New English text。
func (n *Notifier) Run(ctx context.Context) {
	if n.pg == nil {
		return
	}
	t := time.NewTicker(notifyTick)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			n.step(ctx)
		}
	}
}

// step English text：English text，English text。
//
// English text、English text——English text。
// English text tick English text，English text。
func (n *Notifier) step(ctx context.Context) {
	if !n.enabled() {
		return
	}
	if _, _, err := n.pg.FanOutPendingEvents(ctx, notifyFanOutPerTick); err != nil {
		log.Printf("[notify] English text: %v", err)
		return
	}
	channels, err := n.pg.ListNotificationChannels(ctx)
	if err != nil {
		log.Printf("[notify] English text: %v", err)
		return
	}
	baseURL := n.publicBaseURL()
	for _, ch := range channels {
		if !ch.IsEnabled() {
			continue
		}
		// English text**English text**（English text HTTP English text），English text。
		// English text（English text）；English text
		// English text，English text。
		//
		// English text、English text——English text，English text
		// English text。
		now := time.Now()
		if ch.Mode == db.NotifyModeDigest {
			tokens, claimLimit := digestTickPlan()
			if n.takeTokens(ch.ID, ch.RatePerMin, tokens, now) <= 0 {
				continue
			}
			n.stepDigest(ctx, ch, claimLimit, baseURL)
			continue
		}
		allow := n.takeTokens(ch.ID, ch.RatePerMin, notifyMaxSendsPerChannelPerTick, now)
		if allow <= 0 {
			continue
		}
		n.stepRealtime(ctx, ch, allow, baseURL)
	}
}

// digestTickPlan English text。
//
// English text**English text**，English text：
//
//   - tokens English text。English text、English text HTTP English text，English text 1。
//     rate_per_min English text digest English text（English text）。
//   - claimLimit English text。English text，English text。
//
// English text rate_per_min English text digest English text，English text
// （notifyMaxSendsPerChannelPerTick，English text）English text。
// English text rate_per_min=20 English text 3 English text tick English text 1 English text，English text
// English text 1 English text——digest English text「English text」，English text
// 「English text 30 English text 1 English text」，English text db.MaxDigestBatchSize English text。
//
// English text（English text limit English text
// stepDigest，English text step English text），English text
// TestDigestTickPlanDecouplesBatchSizeFromSendBudget English text。
func digestTickPlan() (tokens, claimLimit int) {
	return 1, db.MaxDigestBatchSize
}

// stepRealtime English text，English text。
func (n *Notifier) stepRealtime(ctx context.Context, ch *db.NotificationChannel, allow int, baseURL string) {
	deliveries, err := n.pg.ClaimRealtimeDeliveries(ctx, ch.ID, allow, notifyLease)
	if err != nil {
		log.Printf("[notify] English text channel=%d: %v", ch.ID, err)
		return
	}
	if len(deliveries) == 0 {
		return
	}
	channel, cfg, ok := n.adapt(ch)
	if !ok {
		_ = n.pg.FailDeliveries(ctx, deliveryIDs(deliveries), fmt.Sprintf("English text %q English text", ch.Kind))
		return
	}
	for _, dl := range deliveries {
		msg, err := n.renderSingle(ctx, dl, baseURL)
		if err != nil {
			// English text，English text。
			_ = n.pg.FailDeliveries(ctx, []int64{dl.ID}, err.Error())
			continue
		}
		n.send(ctx, channel, cfg, msg, []*db.NotificationDelivery{dl})
	}
}

// stepDigest English text。
func (n *Notifier) stepDigest(ctx context.Context, ch *db.NotificationChannel, allow int, baseURL string) {
	window := n.digestInterval()
	due, err := n.pg.DigestBatchDue(ctx, ch.ID, window)
	if err != nil {
		log.Printf("[notify] English text channel=%d: %v", ch.ID, err)
		return
	}
	if !due {
		return
	}
	deliveries, err := n.pg.ClaimDigestBatch(ctx, ch.ID, allow, notifyLease)
	if err != nil {
		log.Printf("[notify] English text channel=%d: %v", ch.ID, err)
		return
	}
	if len(deliveries) == 0 {
		return
	}
	channel, cfg, ok := n.adapt(ch)
	if !ok {
		_ = n.pg.FailDeliveries(ctx, deliveryIDs(deliveries), fmt.Sprintf("English text %q English text", ch.Kind))
		return
	}
	msg, included, err := n.renderBatch(ctx, deliveries, baseURL, int(window.Minutes()))
	if err != nil {
		_ = n.pg.FailDeliveries(ctx, deliveryIDs(deliveries), err.Error())
		return
	}
	// English text、English text。English text
	// included English text、English text——English text
	// English text，English text sending English text。
	if skipped := excludeDeliveries(deliveries, included); len(skipped) > 0 {
		reason := "English text，English text"
		if fErr := n.pg.FailDeliveries(ctx, deliveryIDs(skipped), reason); fErr != nil {
			log.Printf("[notify] English text channel=%s ids=%v: %v", ch.Kind, deliveryIDs(skipped), fErr)
		}
		log.Printf("[notify] English text %d English text channel=%d", len(skipped), ch.ID)
	}
	// English text send：included[i] English text msg.Items[i] English text，
	// send English text「English text K English text」English text。
	n.send(ctx, channel, cfg, msg, included)
}

// send English text。
//
// English text（English text）English text：English text、English text。
// English text——English text，English text。
//
// English text**English text**：English text K English text，
// English text K+1 English text，English text。English text
// English text、English text，English text。
func (n *Notifier) send(ctx context.Context, channel notify.Channel, cfg map[string]any, msg notify.Message, deliveries []*db.NotificationDelivery) {
	// English text，English text。
	sendCtx, cancel := context.WithTimeout(ctx, notifySendTimeout)
	defer cancel()
	delivered, err := channel.Send(sendCtx, cfg, msg)
	if err == nil && delivered > 0 {
		if delivered > len(deliveries) {
			// English text；English text，
			// English text，English text。
			log.Printf("[notify] English text %d English text %d channel=%s，English text",
				delivered, len(deliveries), channel.Kind())
			delivered = len(deliveries)
		}
		sent, rest := deliveries[:delivered], deliveries[delivered:]
		if err := n.pg.MarkDeliveriesSent(ctx, deliveryIDs(sent)); err != nil {
			log.Printf("[notify] English text channel=%s ids=%v: %v", channel.Kind(), deliveryIDs(sent), err)
		}
		if len(rest) > 0 {
			// English text：English text，English text tick English text。
			// English text DeferDeliveries English text RescheduleDeliveries —— English text，
			// English text（English text +1 English text，English text）。
			if err := n.pg.DeferDeliveries(ctx, deliveryIDs(rest),
				fmt.Sprintf("English text，English text %d English text，English text", delivered)); err != nil {
				log.Printf("[notify] English text channel=%s ids=%v: %v", channel.Kind(), deliveryIDs(rest), err)
			}
		}
		return
	}
	if err == nil {
		// English text。English text（English text），
		// English text。
		err = fmt.Errorf("English text（delivered=%d）", delivered)
	}

	// English text**English text**English text，English text。
	//
	// English text `if maxAttempts(deliveries) >= MaxNotifyAttempts` English text，English text
	// English text：English text（attempts=2）English text
	// English text（attempts=1）English text failed——English text，
	// English text「English text」English text。
	permanent := notify.IsPermanent(err)
	var failIDs, exhaustedIDs []int64
	byDelay := map[time.Duration][]int64{}
	for _, dl := range deliveries {
		switch {
		case permanent:
			failIDs = append(failIDs, dl.ID)
		case dl.Attempts >= db.MaxNotifyAttempts:
			exhaustedIDs = append(exhaustedIDs, dl.ID)
		default:
			delay := notifyBackoff[min(dl.Attempts, len(notifyBackoff)-1)]
			byDelay[delay] = append(byDelay[delay], dl.ID)
		}
	}

	if len(failIDs) > 0 {
		if fErr := n.pg.FailDeliveries(ctx, failIDs, err.Error()); fErr != nil {
			log.Printf("[notify] English text channel=%s ids=%v: %v", channel.Kind(), failIDs, fErr)
		}
	}
	if len(exhaustedIDs) > 0 {
		reason := fmt.Sprintf("English text %d English text: %s", db.MaxNotifyAttempts, err)
		if fErr := n.pg.FailDeliveries(ctx, exhaustedIDs, reason); fErr != nil {
			log.Printf("[notify] English text channel=%s ids=%v: %v", channel.Kind(), exhaustedIDs, fErr)
		}
	}
	// English text：English text 3 English text，English text，English text
	// UPDATE（English text 500 English text 500 English text）。
	for delay, group := range byDelay {
		if rErr := n.pg.RescheduleDeliveries(ctx, group, delay, err.Error()); rErr != nil {
			log.Printf("[notify] English text channel=%s ids=%v: %v", channel.Kind(), group, rErr)
		}
	}
	if len(failIDs)+len(exhaustedIDs) > 0 {
		log.Printf("[notify] English text channel=%d kind=%s English text=%d English text=%d English text=%d: %s",
			deliveries[0].ChannelID, channel.Kind(), len(failIDs), len(exhaustedIDs), len(byDelay), err)
	}
}

// excludeDeliveries English text all English text keep English text（English text）。
// English text「English text」English text——English text，English text。
func excludeDeliveries(all, keep []*db.NotificationDelivery) []*db.NotificationDelivery {
	inKeep := make(map[*db.NotificationDelivery]bool, len(keep))
	for _, dl := range keep {
		inKeep[dl] = true
	}
	var out []*db.NotificationDelivery
	for _, dl := range all {
		if !inKeep[dl] {
			out = append(out, dl)
		}
	}
	return out
}

// adapt English text。
// English text ok=false English text，English text。
func (n *Notifier) adapt(ch *db.NotificationChannel) (notify.Channel, map[string]any, bool) {
	channel, ok := notify.Get(ch.Kind)
	if !ok {
		return nil, nil, false
	}
	var cfg map[string]any
	if len(ch.Config) > 0 {
		// English text map：English text Validate English text「English text」，
		// English text JSON English text。
		_ = json.Unmarshal(ch.Config, &cfg)
	}
	if cfg == nil {
		cfg = map[string]any{}
	}
	return channel, cfg, true
}

// renderSingle English text。
func (n *Notifier) renderSingle(ctx context.Context, dl *db.NotificationDelivery, baseURL string) (notify.Message, error) {
	snap, err := parseSnapshot(dl)
	if err != nil {
		return notify.Message{}, err
	}
	item, err := n.itemFor(ctx, snap, baseURL)
	if err != nil {
		return notify.Message{}, err
	}
	return notify.Message{Items: []notify.Item{item}, HomeURL: baseURL}, nil
}

// renderBatch English text。English text——English text，
// English text。
//
// English text included English text msg.Items **English text**（English text i English text ↔ English text i English text）。
// English text：English text「English text K English text」English text K English text
// English text。English text included English text，
// English text——English text，English text。
// English text，English text stepDigest。
func (n *Notifier) renderBatch(ctx context.Context, deliveries []*db.NotificationDelivery, baseURL string, windowMinutes int) (notify.Message, []*db.NotificationDelivery, error) {
	items := make([]notify.Item, 0, len(deliveries))
	included := make([]*db.NotificationDelivery, 0, len(deliveries))
	for _, dl := range deliveries {
		snap, err := parseSnapshot(dl)
		if err != nil {
			// English text，English text included——English text
			// （English text，English text「English text」English text）。
			log.Printf("[notify] English text delivery=%d: %v", dl.ID, err)
			continue
		}
		item, err := n.itemFor(ctx, snap, baseURL)
		if err != nil {
			return notify.Message{}, nil, err
		}
		items = append(items, item)
		included = append(included, dl)
	}
	if len(items) == 0 {
		return notify.Message{}, nil, fmt.Errorf("English text %d English text", len(deliveries))
	}
	return notify.Message{
		Items:         items,
		Batch:         true,
		WindowMinutes: windowMinutes,
		HomeURL:       baseURL,
	}, included, nil
}

// itemFor English text，English text。
func (n *Notifier) itemFor(ctx context.Context, snap notify.Snapshot, baseURL string) (notify.Item, error) {
	assets, err := n.pg.NotificationAssetNames(ctx, snap.AssetIDs)
	if err != nil {
		// English text：English text，
		// English text。
		log.Printf("[notify] English text finding=%d: %v", snap.FindingID, err)
	}
	item := notify.Item{
		FindingID:  snap.FindingID,
		Name:       snap.Name,
		VulnClass:  snap.VulnClass,
		Severity:   snap.Severity,
		Summary:    snap.Summary,
		Assets:     assets,
		FromStatus: snap.FromStatus,
		ToStatus:   snap.ToStatus,
	}
	if baseURL != "" {
		// English text web/src/app/(main)/function/findings/detail/page.tsx，
		// English text query English text id English text id。
		item.DetailURL = fmt.Sprintf("%s/function/findings/detail?id=%d", baseURL, snap.FindingID)
	}
	return item, nil
}

// takeTokens English text**English text want English text**English text，English text。
//
// English text = English text（English text HTTP English text）。English text；
// English text，English text 1。
//
// English text，English text。ratePerMin<=0 English text，
// English text，English text。
//
// want English text：English text，English text，
// English text、English text——English text，
// English text「English text」English text。
func (n *Notifier) takeTokens(channelID int64, ratePerMin, want int, now time.Time) int {
	if want <= 0 {
		return 0
	}
	if ratePerMin <= 0 {
		return min(want, notifyUnlimitedBurstPerTick)
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	b := n.buckets[channelID]
	if b == nil {
		b = &notifyBucket{tokens: float64(ratePerMin), lastFill: now}
		n.buckets[channelID] = b
	}
	// English text，English text ratePerMin/60 English text。
	if elapsed := now.Sub(b.lastFill).Seconds(); elapsed > 0 {
		b.tokens = minF(float64(ratePerMin), b.tokens+elapsed*float64(ratePerMin)/60)
		b.lastFill = now
	}
	// English text epsilon English text：English text，English text
	// 0.5 + 0.5 English text 0.9999999999，English text int() English text 0——
	// English text。1e-9 English text，English text。
	take := min(int(b.tokens+1e-9), want)
	if take <= 0 {
		return 0
	}
	b.tokens -= float64(take)
	return take
}

// enabled English text。
func (n *Notifier) enabled() bool {
	return n.pg.GetBool(settingNotifyEnabled, true)
}

// publicBaseURL English text，English text。
func (n *Notifier) publicBaseURL() string {
	v, ok, err := n.pg.GetSetting(settingNotifyPublicBaseURL)
	if err != nil || !ok {
		return ""
	}
	return trimTrailingSlash(v)
}

// digestInterval English text，English text。
func (n *Notifier) digestInterval() time.Duration {
	v, ok, err := n.pg.GetSetting(settingNotifyDigestMinutes)
	if err != nil || !ok {
		return time.Duration(notifyDefaultDigestMinutes) * time.Minute
	}
	m := 0
	if _, err := fmt.Sscanf(v, "%d", &m); err != nil || m <= 0 {
		return time.Duration(notifyDefaultDigestMinutes) * time.Minute
	}
	return time.Duration(m) * time.Minute
}

// parseSnapshot English text。
func parseSnapshot(dl *db.NotificationDelivery) (notify.Snapshot, error) {
	var snap notify.Snapshot
	if len(dl.Snapshot) == 0 {
		return snap, fmt.Errorf("English text %d English text", dl.ID)
	}
	if err := json.Unmarshal(dl.Snapshot, &snap); err != nil {
		return snap, fmt.Errorf("English text %d English text: %w", dl.ID, err)
	}
	if snap.Kind == "" {
		// English text，English text。
		snap.Kind = dl.EventKind
	}
	return snap, nil
}

func deliveryIDs(deliveries []*db.NotificationDelivery) []int64 {
	out := make([]int64, 0, len(deliveries))
	for _, dl := range deliveries {
		out = append(out, dl.ID)
	}
	return out
}

func trimTrailingSlash(s string) string {
	for len(s) > 0 && s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	return s
}

func minF(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
