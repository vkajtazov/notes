// Code generated by go-swagger; DO NOT EDIT.

package notes

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetNoteByIDHandlerFunc turns a function with the right signature into a get note by Id handler
type GetNoteByIDHandlerFunc func(GetNoteByIDParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetNoteByIDHandlerFunc) Handle(params GetNoteByIDParams) middleware.Responder {
	return fn(params)
}

// GetNoteByIDHandler interface for that can handle valid get note by Id params
type GetNoteByIDHandler interface {
	Handle(GetNoteByIDParams) middleware.Responder
}

// NewGetNoteByID creates a new http.Handler for the get note by Id operation
func NewGetNoteByID(ctx *middleware.Context, handler GetNoteByIDHandler) *GetNoteByID {
	return &GetNoteByID{Context: ctx, Handler: handler}
}

/*GetNoteByID swagger:route GET /note/{id} notes getNoteById

Get note by id

*/
type GetNoteByID struct {
	Context *middleware.Context
	Handler GetNoteByIDHandler
}

func (o *GetNoteByID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetNoteByIDParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
