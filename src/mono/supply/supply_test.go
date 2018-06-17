package supply_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "mono/supply"
	"github.com/golang/mock/gomock"
	"github.com/cloudfoundry/libbuildpack"
	"os"
	"errors"
)

//go:generate mockgen -source=supply.go --destination=mocks_test.go --package=supply_test

var _ = Describe("Supply", func() {
	var (
		mockCtrl      *gomock.Controller
		mockManifest  *MockManifest
		mockStager    *MockStager
		mockCommand   *MockCommand
		supplier      Supplier
		actualError   error
		expectedError error
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockManifest = NewMockManifest(mockCtrl)
		mockStager = NewMockStager(mockCtrl)
		mockCommand = NewMockCommand(mockCtrl)

		logger := libbuildpack.NewLogger(os.Stdout)

		actualError = nil
		expectedError = nil

		supplier = Supplier{
			Manifest: mockManifest,
			Stager:   mockStager,
			Command:  mockCommand,
			Log:      logger,
		}
	})

	AfterEach(func() {
		actualError = supplier.Run()

		if actualError != nil && expectedError != nil {
			Expect(actualError).To(Equal(expectedError))
		}

		mockCtrl.Finish()
	})

	It("Should return error", func() {
		expectedError = errors.New("command error")
		mockCommand.EXPECT().Output(gomock.Any(), gomock.Any()).Return("", expectedError)

	})

	It("Should execute command", func() {
		mockCommand.EXPECT().Output(gomock.Any(), gomock.Any()).Return("Mono has been installed.", nil)
	})
})
