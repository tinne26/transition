package miniscene

type OverFlags uint32
const (
	overFlagRestorePlayerOnCam OverFlags = 0b0001
	overFlagUnblockPlayer      OverFlags = 0b0010
)

func (self OverFlags) withUnblockAfter(ticks uint8) OverFlags {
	return (self & 0x00FF_FFFF) | (OverFlags(ticks) << 24)
}

func (self OverFlags) GetUnblockTicks() uint64 {
	return uint64((self & 0xFF00_0000) >> 24)
}

func (self OverFlags) HasToRestorePlayerOnCam() bool {
	return (self & overFlagRestorePlayerOnCam) == overFlagRestorePlayerOnCam
}

func (self OverFlags) HasToUnblockPlayer() bool {
	return (self & overFlagUnblockPlayer) == overFlagUnblockPlayer
}
