package bump

var (
	ResolveRepository = (*bumper).resolveRepository
	ListReleases      = (*bumper).listReleases
	CreateRelease     = (*bumper).createRelease
	CurrentVersion    = (*bumper).currentVersion
	IncrementVersion  = incrementVersion
)

func (b *bumper) SetPrompter(prompter Prompter) {
	b.prompter = prompter
}

func (b *bumper) Repository() string {
	return b.repository
}

func (b *bumper) IsCurrent() bool {
	return b.isCurrent
}

func (b *bumper) IsDraft() bool {
	return b.isDraft
}

func (b *bumper) IsPrerelease() bool {
	return b.isPrerelease
}

func (b *bumper) DiscussionCategory() string {
	return b.discussionCategory
}

func (b *bumper) GenerateNotes() bool {
	return b.generateNotes
}

func (b *bumper) Notes() string {
	return b.notes
}

func (b *bumper) NotesFilename() string {
	return b.notesFilename
}

func (b *bumper) Target() string {
	return b.target
}

func (b *bumper) Title() string {
	return b.title
}

func (b *bumper) AssetFiles() []string {
	return b.assetFiles
}

func (b *bumper) BumpType() BumpType {
	return b.bumpType
}

func (b *bumper) Yes() bool {
	return b.yes
}
