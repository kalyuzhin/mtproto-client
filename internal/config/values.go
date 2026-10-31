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

func (c *Realtime) GetArgon2Memory() int32 {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.cfg.values[argon2Memory].Int32()
}

func (c *Realtime) GetArgon2Threads() int32 {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.cfg.values[argon2Threads].Int32()
}

func (c *Realtime) GetArgon2Time() int32 {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.cfg.values[argon2Time].Int32()
}

func (c *Realtime) GetArgon2SaltLength() int32 {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.cfg.values[argon2SaltLen].Int32()
}

func (c *Realtime) GetArgon2KeyLength() int32 {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.cfg.values[argon2KeyLen].Int32()
}
