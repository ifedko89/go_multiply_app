package int_tests

import (
	"github.com/igor-fedko/go_multiply_app/int_tests/utils"
	"os"
	"testing"
)

var env *utils.TestEnv

func TestMain(m *testing.M) {
	t := &testing.T{} // dummy для TestEnv (без него не построим TestEnv)
	env = utils.NewMongoEnv(t)

	code := m.Run()

	env.Close()
	os.Exit(code)
}
