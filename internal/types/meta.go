package types

import (
	"strconv"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Meta struct {
	Owner     string    `json:"owner"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (m *Meta) ExtractFromLabels(objMeta metav1.ObjectMeta) {
	labels := objMeta.Labels
	if labels == nil {
		return
	}

	m.Owner = labels[labelOwner]
	m.CreatedAt = decodeTime(labels[labelCreateAt])
	m.UpdatedAt = decodeTime(labels[labelUpdateAt])
}

func (m *Meta) InjectToLabels(objMeta *metav1.ObjectMeta) {
	if objMeta.Labels == nil {
		objMeta.Labels = make(map[string]string)
	}
	labels := objMeta.Labels
	labels[labelOwner] = m.Owner
	if _, ok := labels[labelCreateAt]; !ok {
		labels[labelCreateAt] = encodeTime(m.CreatedAt)
	}
	labels[labelUpdateAt] = encodeTime(m.UpdatedAt)
	labels[labelCreatedBy] = labelValueCreateBy
}

func encodeTime(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10)
}

func decodeTime(s string) time.Time {
	milli, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}
	}
	return time.Unix(milli, 0)
}

const (
	labelPrefix = "obs.minimax.com/"

	labelOwner     = labelPrefix + "owner"
	labelCreateAt  = labelPrefix + "create-at"
	labelUpdateAt  = labelPrefix + "update-at"
	labelCreatedBy = labelPrefix + "created-by"

	labelValueCreateBy = "vm-access-proxy"
)
