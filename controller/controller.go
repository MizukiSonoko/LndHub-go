package controller

import (
	"encoding/json"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/MizukiSonoko/lnd-gateway/entity"
	"github.com/MizukiSonoko/lnd-gateway/lightning"
	"github.com/MizukiSonoko/lnd-gateway/repository"
)

var (
	repo           repository.UserRepo
	lnd            lightning.Lnd
	identityPubkey string
)

func init() {
	repo = repository.NewUserRepo()
	// In now, using macOS
	lnd = lightning.NewLnd(
		"localhost:10009",
		os.Getenv("HOME")+"/Library/Application Support/Lnd/tls.cert")
	info, err := lnd.GetInfo()
	if err != nil {
		panic("")
	}
	identityPubkey = info.IdentityPubkey
}

func auth(w http.ResponseWriter, r *http.Request) {
	pUserId := r.Form.Get("userId")
	pPassword := r.Form.Get("password")

	amount, err := strconv.Atoi(pAmount)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(
			ErrorResp{Message: "amount should be number"})
		return
	}
	if amount < 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(
			ErrorResp{Message: "amount should be plus"})
		return
	}

	token := middleware.GenerateToken(nil)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(
		TokenResp{Token: token})
}

func addInvoice(w http.ResponseWriter, r *http.Request) {
	ok := service.VerifyRequest(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(
			ErrorResp{Message: "Unauthorized"})
		return
	}
	user := repo.Get("")

	pMemo := r.Form.Get("memo")
	pAmount := r.Form.Get("amount")

	amount, err := strconv.Atoi(pAmount)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(
			ErrorResp{Message: "amount should be number"})
		return
	}
	if amount < 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(
			ErrorResp{Message: "amount should be plus"})
		return
	}

	resp, err := lnd.AddInvoice(pMemo, uint(amount))
	user.AttachUserInvoice(resp.PaymentRequest)
	err = repo.Update(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(
			ErrorResp{Message: "amount should be number"})
		return
	}

}

func payInvoice(w http.ResponseWriter, r *http.Request) {
	ok := service.VerifyRequest(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(
			ErrorResp{Message: "Unauthorized"})
		return
	}

	user := repo.Get("")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(
			ErrorResp{Message: "MethodNotAllowed"})
		return
	}

	pInvoice := r.Form.Get("invoice")
	pAmount := r.Form.Get("amount")

	amount, err := strconv.Atoi(pAmount)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(
			ErrorResp{Message: "amount should be number"})
		return
	}
	if amount < 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(
			ErrorResp{Message: "amount should be plus"})
		return
	}

	resp, err := lnd.DecodePay(pInvoice)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(
			ErrorResp{Message: "Unauthorized"})
		return
	}
	numSatoshis := resp.NumSatoshis
	if numSatoshis == 0 {
		numSatoshis = uint(amount)
	}

	balance := user.Balance()
	if balance >= numSatoshis+uint(math.Floor(float64(numSatoshis)*0.01)) {
		if identityPubkey == resp.Description {
			payee, err := repo.FindByPaymentHash(resp.PaymentHash)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(
					ErrorResp{Message: "Unauthorized"})
				return
			}

			if user.GetPaymentHashState(resp.PaymentHash) == entity.PAYMENT_HASH_STATE_PAIED {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(
					ErrorResp{Message: "Unauthorized"})
				return
			}

			payee.UpdateBalance(payee.Balance() + numSatoshis)
			err = repo.Update(payee)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(
					ErrorResp{Message: "Unauthorized"})
				return
			}

			user.UpdateBalance(balance - numSatoshis)
			user.AttachTransaction(
				*entity.NewTx(
					time.Now(),
					"paid_invoice",
					numSatoshis+uint(math.Floor(float64(numSatoshis)*0.01)),
					uint(math.Floor(float64(numSatoshis)*0.03)),
					resp.Description,
				))
			err = repo.Update(user)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(
					ErrorResp{Message: "Unauthorized"})
				return
			}

			payee.UpdatePaymentHashState(resp.PaymentHash, entity.PAYMENT_HASH_STATE_PAIED)
			return
		}

	} else {
		client, err := lnd.GetSendPaymentClient()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(
				ErrorResp{Message: "Unauthorized"})
			return
		}

		user.UnlockFounds(pInvoice)
		client.Send(lightning.SendRequest{})
	}

}

func getBtc(w http.ResponseWriter, r *http.Request) {
	ok := service.VerifyRequest(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(
			ErrorResp{Message: "Unauthorized"})
		return
	}
	user := repo.Get("")

	address := user.Getaddress()
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode([]string{
		address,
	})
}

func balance(w http.ResponseWriter, r *http.Request) {
	ok := service.VerifyRequest(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(
			ErrorResp{Message: "Unauthorized"})
		return
	}

	type BalanceResponse struct {
		Balance uint `json:"token"`
	}

	user := repo.Get("")

	_ = json.NewEncoder(w).Encode(BalanceResponse{Balance: user.Balance()})
}

func getTxs(w http.ResponseWriter, r *http.Request) {
	ok := service.VerifyRequest(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(
			ErrorResp{Message: "Unauthorized"})
		return
	}

	type InvoiceResponse struct {
		Invoice string `json:"string"`
	}

	user := repo.Get("")
	user.Invoice()
}

func getUserInvoices(w http.ResponseWriter, r *http.Request) {
	ok := service.VerifyRequest(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(
			ErrorResp{Message: "Unauthorized"})
		return
	}

	type InvoiceResponse struct {
		Invoice string `json:"string"`
	}

	user := repo.Get("")

	_ = json.NewEncoder(w).Encode(InvoiceResponse{Invoice: user.Invoice()})
}

type TokenResp struct {
	Token string `json:"token"`
}

type ErrorResp struct {
	Message string `json:"message"`
}

func GetHandlerFuncs() map[string]func(w http.ResponseWriter, r *http.Request) {
	return map[string]func(w http.ResponseWriter, r *http.Request){}
}
