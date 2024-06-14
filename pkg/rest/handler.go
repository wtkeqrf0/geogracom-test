package rest

import (
	"context"
	"geogracom-test/pkg/kernel"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// API реализует интерфейс [http.Handler] с шаблонным представлением структуры запроса и
// ответа. С помощью этого обработчика можно обернуть функции бизнес логики и не производить
// маршалинг и анмашалинг структур и URL запроса. Например:
//
//	type Server struct {}
//
//	type Request struct {
//		Name string `schema:"name"`
//	}
//
//	type Response struct {
//		Greeting string `json:"greeting"`
//	}
//
//	func (s *Server) Hello(ctx context.Context, r Request) (*Response, error) { ... }
//
//	func main() {
//		var s Server
//		http.Handle("/", API[Request, Response](s.Hello))
//		http.ListenAndServe("localhost:8080", nil)
//	}
//
// Обработчик возвращает статус 200 и JSON ответ с полями:
//
//	struct {
//		Status  string `json:"status"`
//		Message string `json:"message,omitempty"`
//		Payload any    `json:"payload"`
//	}
//
// где:
//   - status  - "success" или "error", в зависимости от того, как завершился вызов обработчика.
//   - message - это err.Error() в случае, когда вызов обработчика завершен с ошибкой.
//   - payload - это JSON вид Response шаблона либо JSON вид error структуры ошибки.
//
// При поступлении запроса в Handler, производится попытка декодирования http тела запроса в
// типизированный Request параметр. Ошибка при этом не проверяется. Далее при помощи библиотеки
// https://github.com/gorilla/schema дополнительно производится попытка применить URL ключи запроса
// на Request. Таким образом, Request по приоритету сначала заполняется из тела запроса, затем
// из URL аргументов.
//
// Далее производится попытка "валидации" запроса. Response проверяется на наличие метода
// Validate() error. Если метод реализован, то он вызывается. Если произошла ошибка валидации,
// то обработчик заполняет поле message и payload и возвращает результат с ошибкой.
//
// При успешной валидации запроса или отсутствия реализации метода Validate у Request выполняется
// функция обработчика. В которую передаются заполненные поля Request.
//
// В зависимости от результата выполнения обработчика заполняются поля структуры и её статус.
// Ответ отправляется в виде JSON со статусом 200.
type API[Request any] func(context.Context, *Request) (any, int, error)

// Handler возвращает [fiber.Handler] для упрощения роутинга в некоторых ситуациях.
func (h API[Request]) Handler(parseFunc func(c *fiber.Ctx, out any) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(Request)

		if err := parseFunc(c, req); err != nil {
			return err
		}

		resp, statusCode, err := h(c.Context(), req)
		if err != nil {
			return err
		}

		c.Status(statusCode)

		if resp != nil {
			return c.JSON(resp)
		}
		return c.JSON(kernel.Response{Description: http.StatusText(statusCode)})
	}
}
