package bump

var (
	ResolveRepository = (*bumper).resolveRepository
	ListReleases      = (*bumper).listReleases
	CreateRelease     = (*bumper).createRelease
)

func (b *bumper) Repository() string {
	return b.repository
}

func (b *bumper) IsCurrent() bool {
	return b.isCurrent
}
