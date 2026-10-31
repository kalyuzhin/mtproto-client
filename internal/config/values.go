package config

import "path"

const (
	sessionFile   configKey = "session_file"
	sessionDir    configKey = "session_directory"
	argon2Memory  configKey = "argon2_memory"
	argon2Threads configKey = "argon2_threads"
	argon2SaltLen configKey = "argon2_salt_length"
	argon2KeyLen  configKey = "argon2_key_length"
	argon2Time    configKey = "argon2_time"

	appID       = "app_id"
	appHash     = "app_hash"
	phoneNumber = "phone"
	password    = "pswd"
)

func (c *Realtime) GetSessionFile() string {
	c.m.RLock()
	defer c.m.RUnlock()

	return path.Join(c.cfg.values[sessionDir].String(),
		c.cfg.values[sessionFile].String())
}

func (c *Realtime) GetAppID() int {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.cfg.appID
}

func (c *Realtime) GetAppHash() string {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.cfg.appHash
}

func (c *Realtime) GetPassword() string {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.cfg.password
}

func (c *Realtime) GetPhoneNumber() string {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.cfg.phoneNumber
}

func (c *Realtime) GetArgon2Memory() uint32 {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.cfg.values[argon2Memory].UInt32()
}

func (c *Realtime) GetArgon2Threads() uint8 {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.cfg.values[argon2Threads].UInt8()
}

func (c *Realtime) GetArgon2Time() uint32 {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.cfg.values[argon2Time].UInt32()
}

func (c *Realtime) GetArgon2SaltLength() uint32 {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.cfg.values[argon2SaltLen].UInt32()
}

func (c *Realtime) GetArgon2KeyLength() uint32 {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.cfg.values[argon2KeyLen].UInt32()
}
