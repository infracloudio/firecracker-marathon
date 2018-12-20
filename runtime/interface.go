package runtime

type UUID string

// type Executor interface {
// 	CreateEnvironment() error
// 	StartEnvironment() error
// 	StopEnvironment() error
// 	AttachCodeToEnvironment() error
// 	ExecuteCodeInEnvironment() error
// 	TerminateEnvironment() error
// }

type Runtime interface {
	CreateEnv(language string) UUID            // Creates microvm with Env, Embed Enviroment Server, Returns UUID
	AttachCodeToEnv(uuid UUID, tarPath string) // Attach code to MicroVM
	ValidateEnv(uuid UUID)                     // Make sure env server loads code, /v2/specialize API call and it should get 200 OK  Make sure everything fine and call is loaded may be by calling v2/specialize api
	CallEnvFunction(uuid UUID, payload string) // Call env function endpoint with payload json
	StartEnv(uuid UUID)                        // Starting Environment
	StopEnv(uuid UUID)                         // Stopping/Killing Env
}
