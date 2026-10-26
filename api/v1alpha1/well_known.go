package v1alpha1

// Finalizer, annotation, and label keys that are part of the respawn API contract.
const (
	// Finalizer is added to GameServer and GameServerBackup objects so the
	// operator can run cleanup before they are removed.
	Finalizer = "respawn.bblaker.com/finalizer"

	// AnnotationWakeRequestedAt is set by the waker, as an RFC 3339 timestamp,
	// when a player tries to join a sleeping server.
	AnnotationWakeRequestedAt = "respawn.bblaker.com/wake-requested-at"

	// AnnotationSkipFinalBackup lets a GameServer be deleted when its final backup fails.
	AnnotationSkipFinalBackup = "respawn.bblaker.com/skip-final-backup"

	// LabelGameServer identifies the GameServer that owns an object.
	LabelGameServer = "respawn.bblaker.com/gameserver"

	// LabelScheduled marks GameServerBackups created by a backup schedule.
	LabelScheduled = "respawn.bblaker.com/scheduled"

	// LabelOrphanedFrom marks a retained PersistentVolumeClaim with the name of
	// the GameServer it belonged to.
	LabelOrphanedFrom = "respawn.bblaker.com/orphaned-from"
)

// Condition types used in GameServer status.
const (
	// ConditionReady is true when players can join right now.
	ConditionReady = "Ready"
	// ConditionAvailable is true when a join attempt will succeed or wake the server.
	ConditionAvailable = "Available"
	// ConditionProgressing is true while a start, stop, restore, or rollout is under way.
	ConditionProgressing = "Progressing"
	// ConditionDegraded is true when something needs attention but the server still works.
	ConditionDegraded = "Degraded"
	// ConditionBackupHealthy is true when the last scheduled backup succeeded recently.
	ConditionBackupHealthy = "BackupHealthy"
)
