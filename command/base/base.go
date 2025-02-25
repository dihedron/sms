package base

import (
	"log/slog"
	"os"
	"runtime"
	"runtime/pprof"
)

type Command struct {
	// Endpoint is the API endpoint.
	Endpoint string `short:"e" long:"endpoint" description:"The API endpoint to use." required:"yes" env:"SMS_ENDPOINT" cfg:"endpoint" default:"https://platform.rdcom.com"`
	// SkipVerifyTLS sets whether to skip TLS verification.
	SkipVerifyTLS bool `short:"S" long:"skip-verify-tls" description:"Whether to skip TLS verification." optional:"yes" hidden:"yes" env:"SMS_SKIP_TLS_VERIFY"`
	// EnableDebug sets whether to enable debug info in API calls.
	EnableDebug bool `short:"D" long:"enable-debug" description:"Whether to enable debug info in API calls." optional:"yes" hidden:"true" env:"SMS_ENABLE_DEBUG"`
	// EnableTrace sets whether to enable trace info in API calls.
	EnableTrace bool `short:"T" long:"enable-trace" description:"Whether to enable trace info in API calls." optional:"yes" hidden:"true" env:"SMS_ENABLE_TRACE"`
	// CPUProfile sets the (optional) path of the file for CPU profiling info.
	CPUProfile *string `short:"C" long:"cpu-profile" description:"The (optional) path where the CPU profiler will store its data." optional:"yes" env:"SMS_CPU_PROFILE"`
	// MemProfile sets the (optional) path of the file for memory profiling info.
	MemProfile *string `short:"M" long:"mem-profile" description:"The (optional) path where the memory profiler will store its data." optional:"yes" env:"SMS_MEM_PROFILE"`
}

type TokenCommand struct {
	Command
	// Token is the authentication token.
	Token *string `short:"t" long:"token" description:"The token to use for authentication." required:"yes" env:"SMS_TOKEN" cfg:"token"`
}

type CredentialsCommand struct {
	Command
	// Username is the username to use in API calls' basic authentication.
	Username string `short:"u" long:"username" description:"The username to use for authentication." required:"yes" env:"SMS_USERNAME" cfg:"username"`
	// Password is the password to use in API calls' basic authentication.
	Password string `short:"p" long:"password" description:"The password to use for authentication." required:"yes" env:"SMS_PASSWORD" cfg:"password"`
}

// ProfileCPU creates a pprof CPU profile of the running application.
func (cmd *Command) ProfileCPU() *Closer {
	var f *os.File
	if cmd.CPUProfile != nil {
		var err error
		f, err = os.Create(*cmd.CPUProfile)
		if err != nil {
			slog.Error("could not create CPU profile", "file", cmd.CPUProfile, "error", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			slog.Error("could not start CPU profiler", "error", err)
		}
	}
	return &Closer{
		file: f,
	}
}

// ProfileMemory creates a memory profile of the running application.
func (cmd *Command) ProfileMemory() {
	if cmd.MemProfile != nil {
		f, err := os.Create(*cmd.MemProfile)
		if err != nil {
			slog.Error("could not create memory profile", "file", cmd.MemProfile, "error", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			slog.Error("could not write memory profile", "error", err)
		}
	}
}

// Closer is used to keep a reference to the file used for profiling.
type Closer struct {
	file *os.File
}

// Close closes the file used for profiling.
func (c *Closer) Close() {
	if c.file != nil {
		pprof.StopCPUProfile()
		c.file.Close()
	}
}

const DefaultDateFormat = "Monday, Jan 02, 2006 at 15:04:05 MST"
