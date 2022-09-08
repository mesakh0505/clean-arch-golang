package cmd

import (
	// "time"

	"github.com/LieAlbertTriAdrian/clean-arch-golang/app/config"

	echo "github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	// "github.com/LieAlbertTriAdrian/clean-arch-golang/ebus"
	todoRESTHandler "github.com/LieAlbertTriAdrian/clean-arch-golang/todo/delivery/rest"

	restMiddleware "github.com/LieAlbertTriAdrian/clean-arch-golang/internal/rest/middleware"
)

var restCommand = &cobra.Command{
	Use:   "rest",
	Short: "Start REST server",
	Run:   restServer,
}

func init() {
	rootCmd.AddCommand(restCommand)
}

func restServer(cmd *cobra.Command, args []string) {
	props := config.LoadForServer(EnvFilePath)
	app := rootConfig.App
	e := echo.New()
	// Middleware execution is Descending, the latest appended middleware on `Use` function will be called first
	e.Use(
		restMiddleware.LogErrorMiddleware(),
		restMiddleware.ErrorMiddleware(),
		restMiddleware.EbusInjectorToRequestContext(eBus),
		restMiddleware.SetRequestContextWithTimeout(app.ContextTimeout),
	)

	e.GET("healthcheck", func(c echo.Context) error {
		// TODO (LieAlbertTriAdrian): Call Healtcheck function here

		// e := ebus.Event{
		// 	Name:        "EBUS_TEST",
		// 	Data:        100,
		// 	OccuredTime: time.Now(),
		// }
		// // This will bring panic for now, because the RabbitMQ and SNS handler is not implemented yet
		// ebus.Publish(c.Request().Context(), e)
		return c.String(200, "Calm down bro, I'm really-really healthy Bro!!!")
	})

	todoRESTHandler.InitTodoHandler(e, todoUsecase)
	logrus.Error(e.Start(props.Address))
}
