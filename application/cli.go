package application

import (
	"encoding/json"
	"fmt"
	"github.com/bdemirpolat/kubecd/models"
	"github.com/bdemirpolat/kubecd/validate"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

const category = "application"

// CLIHandlerInterface contains handler functions
type CLIHandlerInterface interface {
	CreateCLIHandler(c *cli.Context) error
	UpdateCLIHandler(c *cli.Context) error
	GetCLIHandler(c *cli.Context) error
	ListCLIHandler(c *cli.Context) error
	DeleteCLIHandler(c *cli.Context) error
	RegisterCommands(app *cli.App)
}

// CLIHandlerStruct
type CLIHandlerStruct struct {
	service ServiceInterface
}

// NewCLIHandler returns new CLIHandlerStruct / CLIHandlerInterface with service
func NewCLIHandler(service ServiceInterface) CLIHandlerInterface {
	return &CLIHandlerStruct{service: service}
}

// compile time proof
var _ CLIHandlerInterface = &CLIHandlerStruct{}

// CreateCLIHandler validates create request and sends to service
func (h *CLIHandlerStruct) CreateCLIHandler(c *cli.Context) error {
	req := models.ApplicationCreateReq{}
	filePath := c.String("file")
	if filePath == "" {
		fmt.Println("file parameter can not be null")
		return nil
	}

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println(errors.Wrap(err, "error while file opening"))
		return nil
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(errors.Wrap(err, "error while file reading"))
		return nil
	}

	err = json.Unmarshal(data, &req)
	if err != nil {
		fmt.Println(errors.Wrap(err, "error while data unmarshalling"))
		return nil
	}

	err = validate.Validate.Struct(req)
	if err != nil {
		fmt.Println(errors.Wrap(err, "error while file validating"))
		return nil
	}

	res, err := h.service.CreateService(req)
	if err != nil {
		fmt.Printf("application can not created, error: %s\n", err.Error())
		return nil
	}
	fmt.Printf("application created successfully, ID: %d\n", res.ID)
	return nil
}

// UpdateCLIHandler validates update request and sends to service
func (h *CLIHandlerStruct) UpdateCLIHandler(c *cli.Context) error {
	req := models.ApplicationUpdateReq{}
	filePath := c.String("file")
	if filePath == "" {
		fmt.Println("file parameter can not be null")
		return nil
	}

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println(errors.Wrap(err, "error while file opening"))
		return nil
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(errors.Wrap(err, "error while file reading"))
		return nil
	}

	err = json.Unmarshal(data, &req)
	if err != nil {
		fmt.Println(errors.Wrap(err, "error while data unmarshalling"))
		return nil
	}

	if req.ID == 0 {
		fmt.Println("\"id\" can not be 0 or empty for update action")
		return nil
	}

	err = validate.Validate.Struct(req)
	if err != nil {
		fmt.Println(errors.Wrap(err, "error while file validating"))
		return nil
	}

	res, err := h.service.UpdateService(req)
	if err != nil {
		fmt.Printf("application can not updated, error: %s\n", err.Error())
		return nil
	}
	fmt.Printf("application updated successfully, ID: %d\n", res.ID)
	return nil
}

// GetCLIHandler validates get request and sends to service
func (h *CLIHandlerStruct) GetCLIHandler(c *cli.Context) error {
	req := models.ApplicationGetReq{}
	req.ID = c.Uint("id")
	res, err := h.service.GetService(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	pretty, err := json.MarshalIndent(*res.Data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(string(pretty))
	return nil
}

// ListCLIHandler validates list request and sends to service
func (h *CLIHandlerStruct) ListCLIHandler(c *cli.Context) error {
	res, err := h.service.ListService(models.ApplicationListReq{})
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	pretty, err := json.MarshalIndent(*res.Data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(string(pretty))
	return nil
}

// DeleteCLIHandler validates get request and sends to service
func (h *CLIHandlerStruct) DeleteCLIHandler(c *cli.Context) error {
	req := models.ApplicationDeleteReq{}
	req.ID = c.Uint("id")
	_, err := h.service.DeleteService(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	fmt.Println("application deleted")
	return nil
}

// RegisterCommands registers cli commands
func (h *CLIHandlerStruct) RegisterCommands(app *cli.App) {
	createCommand := &cli.Command{
		Name:        "create",
		Usage:       "kubecd application create --file=./myjsonfile.json",
		Description: "create new application",
		Category:    category,
		Action:      h.CreateCLIHandler,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "file"},
		},
	}

	updateCommand := &cli.Command{
		Name:        "update",
		Usage:       "kubecd application update --file=./myjsonfile.json",
		Description: "update application",
		Category:    category,
		Action:      h.UpdateCLIHandler,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "file"},
		},
	}

	getCommand := &cli.Command{
		Name:        "get",
		Usage:       "kubecd application get --id=1",
		Description: "get application details",
		Category:    category,
		Action:      h.GetCLIHandler,
		Flags: []cli.Flag{
			&cli.UintFlag{Name: "id"},
		},
	}

	listCommand := &cli.Command{
		Name:        "list",
		Usage:       "kubecd application list",
		Description: "get all of application details",
		Category:    category,
		Action:      h.ListCLIHandler,
	}

	var applicationSubCommands []*cli.Command
	applicationSubCommands = append(applicationSubCommands, createCommand)
	applicationSubCommands = append(applicationSubCommands, updateCommand)
	applicationSubCommands = append(applicationSubCommands, getCommand)
	applicationSubCommands = append(applicationSubCommands, listCommand)

	applicationCommand := &cli.Command{
		Name:        "application",
		Usage:       "application usage",
		UsageText:   "application usage text",
		Description: "application description",
		Subcommands: applicationSubCommands,
	}

	app.Commands = append(app.Commands, applicationCommand)
}
