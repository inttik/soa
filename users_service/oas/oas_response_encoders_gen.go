// Code generated by ogen, DO NOT EDIT.

package oas

import (
	"io"
	"net/http"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func encodeFriendsUserIDGetResponse(response FriendsUserIDGetRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *FriendList:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *FriendsUserIDGetBadRequest:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		writer := w
		if _, err := io.Copy(writer, response); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *FriendsUserIDGetForbidden:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(403)
		span.SetStatus(codes.Error, http.StatusText(403))

		writer := w
		if _, err := io.Copy(writer, response); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *FriendsUserIDGetNotFound:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(404)
		span.SetStatus(codes.Error, http.StatusText(404))

		writer := w
		if _, err := io.Copy(writer, response); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeFriendsUserIDPostResponse(response FriendsUserIDPostRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *FriendsUserIDPostOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *FriendsUserIDPostBadRequest:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		writer := w
		if _, err := io.Copy(writer, response); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *FriendsUserIDPostUnauthorized:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(401)
		span.SetStatus(codes.Error, http.StatusText(401))

		writer := w
		if _, err := io.Copy(writer, response); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *FriendsUserIDPostForbidden:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(403)
		span.SetStatus(codes.Error, http.StatusText(403))

		writer := w
		if _, err := io.Copy(writer, response); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *FriendsUserIDPostNotFound:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(404)
		span.SetStatus(codes.Error, http.StatusText(404))

		writer := w
		if _, err := io.Copy(writer, response); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeLoginPostResponse(response LoginPostRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *LoginPostOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *LoginPostBadRequest:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		writer := w
		if _, err := io.Copy(writer, response); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *LoginPostNotFound:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(404)
		span.SetStatus(codes.Error, http.StatusText(404))

		writer := w
		if _, err := io.Copy(writer, response); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeProfileUserIDGetResponse(response ProfileUserIDGetRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *ProfileInfo:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *ProfileUserIDGetNotFound:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(404)
		span.SetStatus(codes.Error, http.StatusText(404))

		writer := w
		if _, err := io.Copy(writer, response); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeProfileUserIDPostResponse(response ProfileUserIDPostRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *ProfileInfo:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *ProfileUserIDPostBadRequest:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		writer := w
		if _, err := io.Copy(writer, response); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *ProfileUserIDPostUnauthorized:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(401)
		span.SetStatus(codes.Error, http.StatusText(401))

		writer := w
		if _, err := io.Copy(writer, response); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *ProfileUserIDPostForbidden:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(403)
		span.SetStatus(codes.Error, http.StatusText(403))

		writer := w
		if _, err := io.Copy(writer, response); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *ProfileUserIDPostNotFound:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(404)
		span.SetStatus(codes.Error, http.StatusText(404))

		writer := w
		if _, err := io.Copy(writer, response); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeRegisterPostResponse(response RegisterPostRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *UserId:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(201)
		span.SetStatus(codes.Ok, http.StatusText(201))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *RegisterPostBadRequest:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(400)
		span.SetStatus(codes.Error, http.StatusText(400))

		writer := w
		if _, err := io.Copy(writer, response); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *RegisterPostForbidden:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(403)
		span.SetStatus(codes.Error, http.StatusText(403))

		writer := w
		if _, err := io.Copy(writer, response); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeUserLoginGetResponse(response UserLoginGetRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *UserId:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *UserLoginGetNotFound:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(404)
		span.SetStatus(codes.Error, http.StatusText(404))

		writer := w
		if _, err := io.Copy(writer, response); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeWhoamiGetResponse(response WhoamiGetRes, w http.ResponseWriter, span trace.Span) error {
	switch response := response.(type) {
	case *WhoamiGetOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		span.SetStatus(codes.Ok, http.StatusText(200))

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *WhoamiGetForbidden:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(403)
		span.SetStatus(codes.Error, http.StatusText(403))

		writer := w
		if _, err := io.Copy(writer, response); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}
