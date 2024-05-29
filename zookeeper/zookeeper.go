package zookeeper

import (
	"fmt"
	"ollie/stacks"
	"ollie/styles"
	"sort"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/go-zookeeper/zk"
)

func init() {
	// Disable ZooKeeper client logging by setting the logger to a no-op logger
	zk.DefaultLogger = &noopLogger{}
}

type noopLogger struct{}

func (l *noopLogger) Printf(format string, a ...interface{}) {
	// No-op
}

func (l *noopLogger) Println(v ...interface{}) {
	// No-op
}

func getEnv() string {
	var level string
	form := huh.NewSelect[string]().Title("Select an environment level").
		Options(huh.NewOption("Dev Stack", "stack"),
			huh.NewOption("Staging", "staging"),
		).
		Value(&level)

	form.Run()

	var env string
	var err error
	switch level {
	case "stack":
		env, err = stacks.SelectStack()
		if err != nil {
			log.Fatal("There was an issue getting the stack", err)
		}
	case "staging":
		env = "staging"
	default:
		log.Fatal("Invalid environment level")
	}

	return env
}

func EnterZookeeper() {
	env := getEnv()

	log.Debugf("Searching zookeeper for env %s", env)

	var task string
	form := huh.NewSelect[string]().Title("What do you want to do?").
		Options(
			huh.NewOption("Read What Tracfone is pointing to", "read"),
			huh.NewOption("Set Vendor Environment", "set"),
		).
		Value(&task)

	form.Run()
	switch task {
	case "read":
		readEnv(env)
	case "set":
		setVendorEnv(env)
	default:
		log.Fatal("Invalid task")
	}
}

func readEnv(env string) {
	conn := initializeConnection(env)
	defer conn.Close()

	// Path to read data from
	path := "/vidapay/vendors/tracfone/current_env"

	// Read data from the path
	data, _, err := conn.Get(path)
	if err != nil {
		// Handle errors (e.g., zk.ErrNodeNotFound)
		log.Fatal("Error reading data:", err)
	}

	fmt.Println(styles.HighlightStyle.Render("Current_Env: ", string(data)))
}

func setVendorEnv(env string) {
	conn := initializeConnection(env)
	defer conn.Close()

	data, _, err := conn.Children("/vidapay/vendors")
	if err != nil {
		log.Fatal("Error reading data:", err)
	}
	sort.Strings(data)

	var vendorName string
	vendorForm := huh.NewSelect[string]().Title("Select a vendor?").
		Options(huh.NewOptions(data...)...).
		Value(&vendorName)
	vendorForm.WithAccessible(true).Run()

	currentEnv, _, err := conn.Get(fmt.Sprintf("/vidapay/vendors/%s/current_env", vendorName))
	if err != nil {
		log.Fatal("Error reading current_env:", err)
	}

	fmt.Println(styles.HighlightStyle.Render("Current_Env:", string(currentEnv)))

	var new_env string
	envForm := huh.NewInput().Title("Enter a new environment").
		Value(&new_env)
	envForm.Run()

	_, setError := conn.Set(fmt.Sprintf("/vidapay/vendors/%s/current_env", vendorName), []byte(new_env), -1)
	if setError != nil {
		log.Fatal("Error setting current_env:", setError)
	}

	next_env, _, err := conn.Get(fmt.Sprintf("/vidapay/vendors/%s/current_env", vendorName))
	if err != nil {
		log.Fatal("Error reading current_env:", err)
	}

	fmt.Println(styles.HighlightStyle.Render("Environment now set to:", string(next_env)))
}

func initializeConnection(env string) *zk.Conn {
	// Zookeeper server address
	serverAddr := []string{fmt.Sprintf("zookeeper.%s.tcetra.dev", env)}
	// Connect to Zookeeper
	conn, _, err := zk.Connect(serverAddr, time.Second*10)
	if err != nil {
		log.Fatalf("Error connecting to ZooKeeper: %s", err)
	}

	return conn
}
