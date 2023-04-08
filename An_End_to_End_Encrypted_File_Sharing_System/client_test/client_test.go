package client_test

// You MUST NOT change these default imports.  ANY additional imports it will
// break the autograder and everyone will be sad.

import (
	// Some imports use an underscore to prevent the compiler from complaining
	// about unused imports. Normally, you will want to avoid underscore imports
	// unless you know exactly what you are doing. You can read more about
	// underscore imports here: https://golangdocs.com/blank-identifier-in-golang
	_ "encoding/hex"
	_ "errors"
	_ "strconv"
	_ "strings"
	"testing"

	// A "dot" import is used here so that the functions in the ginko and gomega
	// modules can be used without an identifier. For example, Describe() and
	// Expect() instead of ginko.Describe() and gomega.Expect(). You can read more
	// about dot imports here:
	// https://stackoverflow.com/questions/6478962/what-does-the-dot-or-period-in-a-go-import-statement-do
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	userlib "github.com/cs161-staff/project2-userlib"

	// The client implementation is intentionally defined in a different package.
	// This forces us to follow best practice and write tests that only rely on
	// client API that is exported from the client package, and avoid relying on
	// implementation details private to the client package.
	"github.com/cs161-staff/project2-starter-code/client"
)

func TestSetupAndExecution(t *testing.T) {
	// We are using 2 libraries to help us write readable and maintainable tests:
	//
	// (1) Ginkgo, a Behavior Driven Development (BDD) testing framework that
	//             makes it easy to write expressive specs that describe the
	//             behavior of your code in an organized manner; and
	//
	// (2) Gomega, an assertion/matcher library that allows us to write individual
	//             assertion statements in tests that read more like natural
	//             language. For example "Expect(ACTUAL).To(Equal(EXPECTED))".
	//
	// In the Ginko framework, a test case signals failure by calling Ginkgo’s
	// Fail(description string) function. However, we are using the Gomega library
	// to execute our assertion statements. When a Gomega assertion fails, Gomega
	// calls a GomegaFailHandler, which is a function that must be provided using
	// gomega.RegisterFailHandler(). Here, we pass Ginko's Fail() function to
	// Gomega so that Gomega can report failed assertions to the Ginko test
	// framework, which can take the appropriate action when a test fails.
	//
	// This is the sole connection point between Ginkgo and Gomega.
	RegisterFailHandler(Fail)

	RunSpecs(t, "Client Tests")
}

// ================================================
// Here are some optional global variables that can be used throughout the test
// suite to make the tests more readable and maintainable than defining these
// values in each test. You can add more variables here if you want and think
// they will help keep your code clean!
// ================================================
const someFilename = "file1.txt"
const someOtherFilename = "file2.txt"
const nonExistentFilename = "thisFileDoesNotExist.txt"

const aliceUsername = "Alice"
const alicePassword = "AlicePassword"
const bobUsername = "Bob"
const bobPassword = "BobPassword"
const nilufarUsername = "Nilufar"
const nilufarPassword = "NilufarPassword"
const olgaUsername = "Olga"
const olgaPassword = "OlgaPassword"
const marcoUsername = "Marco"
const marcoPassword = "MarcoPassword"

const nonExistentUsername = "NonExistentUser"

var alice *client.User
var bob *client.User
var nilufar *client.User
var olga *client.User
var marco *client.User

var someFileContent []byte
var someShortFileContent []byte
var someLongFileContent []byte

// ================================================
// The top level Describe() contains all tests in
// this test suite in nested Describe() blocks.
// ================================================

var _ = Describe("Client Tests", func() {
	BeforeEach(func() {
		// This top-level BeforeEach will be run before each test.
		//
		// Resets the state of Datastore and Keystore so that tests do not
		// interfere with each other.
		userlib.DatastoreClear()
		userlib.KeystoreClear()

		userlib.SymbolicDebug = false
		userlib.SymbolicVerbose = false
	})

	BeforeEach(func() {
		// This top-level BeforeEach will be run before each test.
		//
		// Byte slices cannot be constant, so this BeforeEach resets the content of
		// each global variable to a predefined value, which allows tests to rely on
		// the expected value of these variables.
		someShortFileContent = []byte("some short file content")
		someFileContent = someShortFileContent
		someLongFileContent = []byte("some LOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOONG file content")
	})

	Describe("Creating users", func() {
		BeforeEach(func() {
			// This top-level BeforeEach will be run before each test.
			//
			// Resets the state of Datastore and Keystore so that tests do not
			// interfere with each other.
			userlib.DatastoreClear()
			userlib.KeystoreClear()

			userlib.SymbolicDebug = false
			userlib.SymbolicVerbose = false
		})

		It("should not error when creating a new user", func() {
			_, err := client.InitUser("Alice", "password")
			Expect(err).To(BeNil(), "Failed to initialize user Alice.")
		})

		It("should error if a username is already taken by another user", func() {
			// TODO: implement this function
			alice1, err1 := client.InitUser("user_alice", "password1")
			alice2, err2 := client.InitUser("user_alice", "password2")

			if alice1 != nil || err1 == nil {
				Expect(err1).To(BeNil(), "Failed to make err nill.")
				Expect(alice1).ToNot(BeNil(), "Failed to initialize user alice1 correctly.")
			}

			if err2 != nil || alice2 == nil {
				Expect(alice2).To(BeNil(), "Failed to create unique usernames.")
				Expect(err2).ToNot(BeNil(), "Failed to make err nill.")
			}

			//fmt.Println("My name is Delina and I am so over life rn")
		})

		It("should error if a user does not exist with that username", func() {
			// TODO: implement function
			bob1, _ := client.GetUser("user_bob", "password2")
			Expect(bob1).To(BeNil(), "Failed to error when a user does not exist with that username.")
		})

		// TODO: you probably want more test cases about creating users here
		It("should error if a user gives the wrong password when logging in", func() {
			// TODO: implement function
			mike1, _ := client.InitUser("user_mike", "password3")
			oh_no, err := client.GetUser("user_mike", "password2")

			Expect(oh_no).To(BeNil(), "Failed to correctly check/return if user put in correct password when logging in.")
			Expect(err).ToNot(BeNil(), "Failed to throw error when user logs in with incorrect password.")
			Expect(mike1.Username).To(BeEquivalentTo("user_mike"), "Failed to return correct user. Username should me: user_mike.")

			/*
				username_in_bytes := []byte(mike1.Username)
				password_in_bytes := []byte("password2")
				testing_password := userlib.Argon2Key(password_in_bytes, username_in_bytes, 48)

				Expect(testing_password).ToNot(BeEquivalentTo(mike1.Password), "Code used to check if passwords match are incorrect")

				password_in_bytes = []byte("password3")
				testing_password = userlib.Argon2Key(password_in_bytes, username_in_bytes, 48)

				Expect(testing_password).To(BeEquivalentTo(mike1.Password), "Code used to check if passwords match are incorrect.")
			*/

		})

		It("should return the same User struct when loggin in as when created", func() {
			toby1, _ := client.InitUser("user_toby", "password4")
			toby2, _ := client.GetUser("user_toby", "password4")

			Expect(toby1).To(BeEquivalentTo(toby2), "Failed to return correct User Struct. Should be the same as the one created initially.")

		})

		It("should error when new user's username is empty", func() {
			empty_user, _ := client.InitUser("", "password")
			Expect(empty_user).To(BeNil(), "Failed to check if username is empty when creating user.")

			empty_user, _ = client.GetUser("", "password")
			Expect(empty_user).To(BeNil(), "Should not return a user struct if username is empty.")
		})

		It("usernames should be case sensitive", func() {
			_, e1 := client.InitUser("Bob", "password")
			Expect(e1).To(BeNil(), "Failed to check if username is case sensitive.")

			_, e2 := client.InitUser("bob", "password")
			Expect(e2).To(BeNil(), "Should not return a username is case sensitive.")
		})

		It("should work for empty passwords", func() {
			_, e1 := client.InitUser("Bob", "")
			Expect(e1).To(BeNil(), "Failed to check if password could be an empty string.")
		})

		It("should work for empty file names", func() {
			bob, e1 := client.InitUser("user_bob", "password2")
			Expect(e1).To(BeNil(), "Failed to create User.")

			e2 := bob.StoreFile("", []byte("Hello World."))
			Expect(e2).To(BeNil(), "Failed to create File with empty string.")
		})

		/*Before executing any of your tests, make sure to set userlib.SymbolicDebug to true
		(this is set by default), and once you've run what you want, use userlib.DebugExportDatastore
		to save the datastore's state for inspection.*/
		It("should create 2 private keys after a user is created", func() {
			userlib.DatastoreClear()
			userlib.KeystoreClear()

			userlib.SymbolicDebug = true
			_, e1 := client.InitUser("user_bob", "password2")
			Expect(e1).To(BeNil(), "Failed to create User.")

			keystore := userlib.KeystoreGetMap()
			Expect(len(keystore)).To(BeEquivalentTo(2),
				"2 keys should be made when a user is created",
				len(keystore),
				2,
			)
		})

		//The client MUST ensure the integrity of filenames.
		It("should ensure the integrity of filenames", func() {
			userlib.DatastoreClear()
			userlib.KeystoreClear()

			userlib.SymbolicDebug = true
			bob, e1 := client.InitUser("user_bob", "password2")
			Expect(e1).To(BeNil(), "Failed to create User.")

			bob.StoreFile("foo.txt", []byte("Hello world."))

			Uuid_of_user, _ := userlib.UUIDFromBytes(userlib.Hash([]byte("user_bob")))
			datastore := userlib.DatastoreGetMap()
			datastore[Uuid_of_user] = []byte("gotcha")

			bob, e3 := client.GetUser("user_bob", "password2")
			Expect(e3).ToNot(BeNil(), "Integrity of entire user is compromised. Should throw an error.")

			file_downloaded, e2 := bob.LoadFile("foo.txt")
			Expect(e2).ToNot(BeNil(), "Integrity of entire user is compromised. Should throw an error.")
			Expect([]byte("Hello world.")).ToNot(BeEquivalentTo(file_downloaded),
				"Downloaded content is not the same as uploaded content",
				[]byte("Hello world."),
				file_downloaded,
			)
		})

	})

	Describe("Single user storage", func() {
		BeforeEach(func() {
			// This top-level BeforeEach will be run before each test.
			//
			// Resets the state of Datastore and Keystore so that tests do not
			// interfere with each other.
			userlib.DatastoreClear()
			userlib.KeystoreClear()

			userlib.SymbolicDebug = false
			userlib.SymbolicVerbose = false
		})

		var alice *client.User

		BeforeEach(func() {
			// This BeforeEach will run before each test in this Describe block.
			alice, _ = client.InitUser("Alice", "some password")
		})

		It("should upload content without erroring", func() {
			content := []byte("This is a test")
			err := alice.StoreFile("file1", content)
			Expect(err).To(BeNil(), "Failed to upload content to a file", err)
		})

		It("should download the expected content that was previously uploaded", func() {
			uploadedContent := []byte("This is a test")
			alice.StoreFile(someFilename, uploadedContent)
			downloadedContent, _ := alice.LoadFile(someFilename)
			Expect(downloadedContent).To(BeEquivalentTo(uploadedContent),
				"Downloaded content is not the same as uploaded content",
				downloadedContent,
				uploadedContent)
		})

		It("should error when trying to download a file that does not exist", func() {
			_, err := alice.LoadFile(nonExistentFilename)
			Expect(err).ToNot(BeNil(), "Was able to load a non-existent file without error.")
		})

		// TODO: you probably want more test cases for store/load/append with a
		// 			 single user here
		//1. Append to a file
		It("should append to a file", func() {
			content := []byte("This is a ")
			err := alice.StoreFile("file2", content)
			Expect(err).To(BeNil(), "Failed to upload content to a file", err)

			err2 := alice.AppendToFile("file2", []byte("test."))
			Expect(err2).To(BeNil(), "Error not nil-> Something is wrong in AppendToFile() func.", err)

			file_content, err3 := alice.LoadFile("file2")
			Expect(err3).To(BeNil(), "Error not nil-> Something is wrong in LoadFile() func.", err)

			Expect(string(file_content)).To(BeEquivalentTo("This is a test."))

		})

		//2. Should be able to store, load, and append to an empty string
		It("should be able to store & load an empty string/content into a file", func() {
			empty_content := []byte("")
			err := alice.StoreFile("file3", empty_content)
			Expect(err).To(BeNil(), "Failed to store empty string/content to a file", err)

			downloadedContent, _ := alice.LoadFile("file3")
			Expect(downloadedContent).To(BeEquivalentTo(empty_content),
				"Failed to load empty string/content to a file",
				downloadedContent,
				empty_content)

			err2 := alice.AppendToFile("file3", []byte("This is another test."))
			Expect(err2).To(BeNil(), "Error not nil-> Something is wrong in AppendToFile() func.", err)

			file_content, err3 := alice.LoadFile("file3")
			Expect(err3).To(BeNil(), "Error not nil-> Something is wrong in LoadFile() func.", err)

			Expect(string(file_content)).To(BeEquivalentTo("This is another test."))

		})

		//4. Make sure contents are overwritten if there is a file already in namespace
		It("should overwrite contents if filename is already in the user's namespace", func() {
			content := []byte("This is a test")
			err := alice.StoreFile("file4", content)
			Expect(err).To(BeNil(), "Failed to upload content to a file", err)

			overwritten_content := []byte("Testing, testing.. 1, 2, 3")
			err2 := alice.StoreFile("file4", overwritten_content)
			Expect(err2).To(BeNil(), "Failed to upload overwritten content to a file", err)

			downloadedContent, _ := alice.LoadFile("file4")
			Expect(downloadedContent).To(BeEquivalentTo(overwritten_content),
				"Content in file should be overwritten",
				downloadedContent,
				overwritten_content)

		})

		//5. Append to someLongFileContent and see that it appended quickly
		It("should append to a file", func() {
			err := alice.StoreFile("file5", someLongFileContent)
			Expect(err).To(BeNil(), "Failed to upload content to a file", err)

			appending_content := []byte("... Testing: Is it fast?")
			err2 := alice.AppendToFile("file5", appending_content)
			Expect(err2).To(BeNil(), "Error not nil-> Something is wrong in AppendToFile() func.", err)

			file_content, err3 := alice.LoadFile("file5")
			Expect(err3).To(BeNil(), "Error not nil-> Something is wrong in LoadFile() func.", err)

			expected_content := append(someLongFileContent, appending_content...)
			Expect(string(file_content)).To(BeEquivalentTo(string(expected_content)))

		})

		//-> IMLEMENT THE "HELPFUL EXAMPLES AS TEST CASES!!!"
		It("should allow different users to use the application at the same time AND support a single user having multiple active sessions at the same time", func() {
			//Multiple Users using software at the same time
			f1 := []byte("content")
			f2 := []byte("different content")

			// Dina and Eric each start a users session by authenticating to the client.
			dina_session_1, _ := client.InitUser("user_dina", "password1")
			eric_session_1, _ := client.InitUser("user_eric", "password2")

			// Dina stores byte slice f1 with name "filename" and Eric stores byte slice
			// f2 also with name "filename".
			dina_session_1.StoreFile("filename", f1)
			eric_session_1.StoreFile("filename", f2)

			//Dina and Eric each confirm that they can load the file they previously
			// stored and that the file contents is the same.

			f1_loaded, _ := dina_session_1.LoadFile("filename")
			f2_loaded, _ := eric_session_1.LoadFile("filename")

			Expect(f1).To(BeEquivalentTo(f1_loaded),
				"file contents should be the same.",
				f1,
				f1_loaded)
			Expect(f2).To(BeEquivalentTo(f2_loaded),
				"file contents should be the same.",
				f2,
				f2_loaded)
			//Dina and Eric have DIFFERENT files
			Expect(f1).ToNot(BeEquivalentTo(f2),
				"file contents should not be the same.",
				f1,
				f2)

			//Single User with Multiple Sessions
			// Eric creates a second user session by authenticating to the client again.
			eric_session_2, _ := client.GetUser("user_eric", "password2")

			// Eric stores byte slice f2 with name "newfile" using his second user
			// session.
			eric_session_2.StoreFile("newfile", f2)

			// Eric loads "newfile" using his first user session. Notice that Eric does
			// not need to reauthenticate. File changes must be available to all active
			// sessions for a given user.

			f2_newfile, _ := eric_session_1.LoadFile("newfile")

			Expect(f2).To(BeEquivalentTo(f2_newfile),
				"file contents should be the same.",
				f2,
				f2_newfile)

			//Eric session1 appends to the new file Eric session 2 made. New append should show to Eric's session2
			adding_content := []byte("Testing, 1, 2, 3.")
			e2 := eric_session_1.AppendToFile("newfile", adding_content)
			Expect(e2).To(BeNil(), "Appending to existing file (and syncing) has an error.", e2)

			appended_file, e3 := eric_session_2.LoadFile("newfile")
			real_file := string(f2) + "Testing, 1, 2, 3."
			Expect(e3).To(BeNil(), "Loading to file (and syncing) has an error.", e3)
			Expect(appended_file).To(BeEquivalentTo(real_file),
				"file contents should be the same.",
				appended_file,
				real_file)
		})

		//File names are not globally unique
		It("should share a file without erroring", func() {
			/*
				For example, user bob can have a file named foo.txt and user alice can have a file named foo.txt.
				There is a single instance of the Datastore server that is used by all users, but the client application MUST keep
				each user’s file namespace independent from one another.
			*/
			bob, _ := client.InitUser("user_bob", "password2")

			err1 := alice.StoreFile("foo.txt", []byte("some_content"))
			err2 := bob.StoreFile("foo.txt", []byte("some_other_content"))
			Expect(err1).To(BeNil(), "Alice could not store file correctly.")
			Expect(err2).To(BeNil(), "Bob could not store file correctly.")

			alices_content, _ := alice.LoadFile("foo.txt")
			bobs_content, _ := bob.LoadFile("foo.txt")
			Expect(alices_content).ToNot(BeEquivalentTo(bobs_content),
				"The file contents that Bob downloaded should not be the same as what Alice downloaded.")
		})

	})

	Describe("Sharing files", func() {
		BeforeEach(func() {
			// This top-level BeforeEach will be run before each test.
			//
			// Resets the state of Datastore and Keystore so that tests do not
			// interfere with each other.
			userlib.DatastoreClear()
			userlib.KeystoreClear()

			userlib.SymbolicDebug = false
			userlib.SymbolicVerbose = false
		})

		BeforeEach(func() {
			// Initialize each user to ensure the variable has the expected value for
			// the tests in this Describe() block.
			alice, _ = client.InitUser(aliceUsername, alicePassword)
			bob, _ = client.InitUser(bobUsername, bobPassword)
			nilufar, _ = client.InitUser(nilufarUsername, nilufarPassword)
			olga, _ = client.InitUser(olgaUsername, olgaPassword)
			marco, _ = client.InitUser(marcoUsername, marcoPassword)
		})

		It("should share a file without erroring", func() {
			alice.StoreFile(someFilename, someShortFileContent)
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Bob.")

			err = bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not receive the file that Alice shared.")

			downloadedContent, err := bob.LoadFile(someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not load the file that Alice shared.")
			Expect(downloadedContent).To(BeEquivalentTo(someShortFileContent),
				"The file contents that Bob downloaded was not the same as what Alice uploaded.")
		})

		// TODO: you probably want more test cases for sharing files here
		It("should have shared users work on the same shared file.", func() {
			alice.StoreFile(someFilename, someShortFileContent)
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Bob.")

			err = bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not receive the file that Alice shared.")

			//Appending test:
			addthis_content := []byte("... Testing: 1, 2, 3 :)")
			err2 := alice.AppendToFile(someFilename, addthis_content)
			Expect(err2).To(BeNil(), "Cannot append to shared file.", err)

			appended_content := []byte(string(someShortFileContent) + "... Testing: 1, 2, 3 :)")
			downloadedContent, err := bob.LoadFile(someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not load the file that Alice shared.")
			Expect(downloadedContent).To(BeEquivalentTo(appended_content),
				"Appending to a shared file should show on both user's structs.")

			//Overwriting test:
			overwriting_file_contents := []byte("Testing: 1, 2, 3 ... Is this good? :)")
			alice.StoreFile(someFilename, overwriting_file_contents)

			downloadedContent, err = bob.LoadFile(someOtherFilename)
			Expect(err).To(BeNil(), "Bob could not load the file that Alice shared.")
			Expect(downloadedContent).To(BeEquivalentTo(overwriting_file_contents),
				"Overwriting a shared file should show on both user's structs.")
		})

		It("should have tree-like structure.", func() {
			//Alice -> Bob -> Olga
			// Bob -> Marco
			//Alice -> Nullifer
			alice.StoreFile(someFilename, someShortFileContent)

			//A-> B
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Bob.")
			err1 := bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err1).To(BeNil(), "Bob could not accept Alice's shared invitation.")

			//A -> N
			shareFileInfoPtr2, err := alice.CreateInvitation(someFilename, nilufarUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Nilufar.")
			err2 := nilufar.AcceptInvitation(aliceUsername, shareFileInfoPtr2, someOtherFilename)
			Expect(err2).To(BeNil(), "Nilufar could not accept Alice's shared invitation.")

			//B -> O
			shareFileInfoPtr3, err := bob.CreateInvitation(someOtherFilename, olgaUsername)
			Expect(err).To(BeNil(), "Bob failed to share a file with Olga.")
			err3 := olga.AcceptInvitation(bobUsername, shareFileInfoPtr3, someOtherFilename)
			Expect(err3).To(BeNil(), "Olga could not accept Bob's shared invitation.")

			//B -> Marco
			shareFileInfoPtr4, err := bob.CreateInvitation(someOtherFilename, marcoUsername)
			Expect(err).To(BeNil(), "Bob failed to share a file with Marco.")
			err4 := marco.AcceptInvitation(bobUsername, shareFileInfoPtr4, someOtherFilename)
			Expect(err4).To(BeNil(), "Marco could not accept Bob's shared invitation.")

			//If Olga makes a change to file contents, everyone should see the changes
			//Appending test:
			addthis_content := []byte("... Testing: 1, 2, 3 :)")
			err2 = olga.AppendToFile(someOtherFilename, addthis_content)
			Expect(err2).To(BeNil(), "Cannot append to shared file.", err)

			appended_content := []byte(string(someShortFileContent) + "... Testing: 1, 2, 3 :)")
			downloadedContent, err := alice.LoadFile(someFilename)
			Expect(err).To(BeNil(), "Alice (root of tree) could not load the shared file that Olga edited.")
			Expect(downloadedContent).To(BeEquivalentTo(appended_content),
				"Appending to a shared file should show on all user's structs.")
		})

		/* The client MUST enforce authorization for all files. The only users who are authorized to access a file using the client
		include: (1) the owner of the file; and (2) users who have accepted an invitation to access the file and that access has not
		been revoked.
		*/
		It("should enfore authorization for all files", func() {
			//Alice -> Bob
			//Alice -> Nullifer
			alice.StoreFile(someFilename, someShortFileContent)

			//A-> B
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Bob.")
			err1 := bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err1).To(BeNil(), "Bob could not accept Alice's shared invitation.")

			//A -> N
			shareFileInfoPtr2, err := alice.CreateInvitation(someFilename, nilufarUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Nilufar.")
			err2 := nilufar.AcceptInvitation(aliceUsername, shareFileInfoPtr2, someOtherFilename)
			Expect(err2).To(BeNil(), "Nilufar could not accept Alice's shared invitation.")

			//File was never shared to olga... lets see if olga can access shared file
			_, err3 := olga.LoadFile(someOtherFilename)
			Expect(err3).ToNot(BeNil(), "Olga should not be able to access shared file without invitation.")

		})

		It("should make sure users who intercept shared files cannot access the file", func() {
			alice.StoreFile(someFilename, someShortFileContent)

			//Alice shares file with Bob
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Bob.")

			//Nilufar intercepts the file
			err1 := nilufar.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err1).ToNot(BeNil(), "Nilufar should not be able to access/accept file.")

			//Only Bob should be able to access/accept file
			err2 := bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err2).To(BeNil(), "Bob should be able to accept Alice's shared invitation.")
		})

		/* The client MUST enforce authorization for all files. The only users who are authorized to access a file using the
		client include: (1) the owner of the file; and (2) users who have accepted an invitation to access the file and that access
		has not been revoked.
		*/
		It("should make sure users who HAVE ACCEPTED THE INVITATION can access the file", func() {
			alice.StoreFile(someFilename, someShortFileContent)

			//Alice shares file with Bob
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Bob.")

			//Before Bob accepts invitation, he tries to access file -> SHOULD FAIL
			_, e1 := bob.LoadFile(someOtherFilename)
			Expect(e1).ToNot(BeNil(), "User should NOT be able to Load to file until invitation is accepted.")

			e3 := bob.AppendToFile(someOtherFilename, []byte("some more stuff"))
			Expect(e3).ToNot(BeNil(), "User should NOT be able to Append to file until invitation is accepted..")

			//NOW Bob accepts invitation -> should now be able to access file
			err2 := bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err2).To(BeNil(), "Bob should be able to accept Alice's shared invitation.")

			_, e4 := bob.LoadFile(someOtherFilename)
			Expect(e4).To(BeNil(), "User should be able to Load to file after invitation is accepted.")

			e5 := bob.StoreFile(someOtherFilename, []byte("Testing, uno, dos, tres."))
			Expect(e5).To(BeNil(), "User should be able to Store to file after invitation is accepted.")

			e6 := bob.AppendToFile(someOtherFilename, []byte("some more stuff"))
			Expect(e6).To(BeNil(), "User should be able to Append to file after invitation is accepted.")
		})

		//Should ensure INTEGRITY of file contents
		It("should ensure integrity of file contents", func() {
			alice.StoreFile(someFilename, []byte("some_content"))
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Bob.")

			//Malicious user tries to change file (compromising the integrity of the file)
			datastore := userlib.DatastoreGetMap()
			datastore[shareFileInfoPtr] = []byte("gotcha")

			err1 := bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err1).ToNot(BeNil(), "File losing its integrity should cause error.")

		})
	})

	// TODO: you probably want more Describe() blocks to contain tests related to
	//       logical test groupings other than the ones suggested above
	Describe("Revoking User Access", func() {
		BeforeEach(func() {
			// This top-level BeforeEach will be run before each test.
			//
			// Resets the state of Datastore and Keystore so that tests do not
			// interfere with each other.
			userlib.DatastoreClear()
			userlib.KeystoreClear()

			userlib.SymbolicDebug = false
			userlib.SymbolicVerbose = false
		})

		BeforeEach(func() {
			// Initialize each user to ensure the variable has the expected value for
			// the tests in this Describe() block.
			alice, _ = client.InitUser(aliceUsername, alicePassword)
			bob, _ = client.InitUser(bobUsername, bobPassword)
			nilufar, _ = client.InitUser(nilufarUsername, nilufarPassword)
			olga, _ = client.InitUser(olgaUsername, olgaPassword)
			marco, _ = client.InitUser(marcoUsername, marcoPassword)
		})

		It("should pass helpful example", func() {
			//1. Alice owns a file named "file.txt" and shares the file with Bob.
			alice.StoreFile("file.txt", []byte("content"))
			invitationPtr, _ := alice.CreateInvitation("file.txt", bobUsername)
			bob.AcceptInvitation(aliceUsername, invitationPtr, "from_alice.txt")

			//2. Alice also shares the file with Nilufar, and nilufar accepts invite
			invitationPtr2, _ := alice.CreateInvitation("file.txt", nilufarUsername)
			nilufar.AcceptInvitation(aliceUsername, invitationPtr2, "from_alice.txt")

			//3. Next, Alice revokes Bob’s permission to this file:
			e1 := alice.RevokeAccess("file.txt", bobUsername)
			Expect(e1).To(BeNil(), "Alice did not revoke access to Bob correctly.")

			//Now, Bob is no longer authorized to access the file. Alice is always authorized to access the file
			//because Alice is the owner. Nulifer is also allowed to access the file because Alice (the owner) shared it with her; not Bob.
			//The client must ensure that Bob cannot use the client API to access the file.
			//The client must ensure that Nulifar can use the client API to access the file.

			//1. Testing Owner Access
			//Overwritten + Loading Test
			_, e2 := alice.LoadFile("file.txt")
			Expect(e2).To(BeNil(), "Owner should still be able to Load file.")

			e3 := alice.StoreFile("file.txt", []byte("new stuff"))
			alice_data, _ := alice.LoadFile("file.txt")
			Expect(e3).To(BeNil(), "Owner should still be able to Store file.")
			Expect("new stuff").To(BeEquivalentTo(string(alice_data)),
				"File should be overwritten after StoreFile() is called.")

			//Append test
			e4 := alice.AppendToFile("file.txt", []byte("some more stuff"))
			Expect(e4).To(BeNil(), "Owner should still be able to Append to file.")

			alice_data2, _ := alice.LoadFile("file.txt")
			real_data := "new stuff" + "some more stuff"
			Expect(real_data).To(BeEquivalentTo(string(alice_data2)),
				"File should be the same.")

			//2. Testing Revoke Access
			_, e5 := bob.LoadFile("from_alice.txt")
			Expect(e5).ToNot(BeNil(), "Revoked User should NOT be able to Load to file.")

			e6 := bob.StoreFile("from_alice.txt", []byte("Testing, uno, dos, tres."))
			Expect(e6).ToNot(BeNil(), "Revoked User should NOT be able to Store to file.")

			e7 := bob.AppendToFile("from_alice.txt", []byte("some more stuff"))
			Expect(e7).ToNot(BeNil(), "Revoked User should NOT be able to Append to file.")

			//3. Testing User who should still have access
			//Overwritten + Loading Test
			e3 = nilufar.StoreFile("from_alice.txt", []byte("new stuff"))
			Expect(e3).To(BeNil(), "Owner should still be able to Store file.")

			nilufar_data, e2 := nilufar.LoadFile("from_alice.txt")
			Expect(e2).To(BeNil(), "Owner should still be able to Load file.")
			Expect("new stuff").To(BeEquivalentTo(string(nilufar_data)),
				"File should be overwritten after StoreFile() is called.")

			//Append test
			e4 = nilufar.AppendToFile("from_alice.txt", []byte("some more stuff"))
			Expect(e4).To(BeNil(), "Owner should still be able to Append to file.")

			nilufar_data2, _ := nilufar.LoadFile("from_alice.txt")
			real_data = "new stuff" + "some more stuff"
			Expect(real_data).To(BeEquivalentTo(string(nilufar_data2)),
				"File should be the same.")
		})

		It("should make sure only the owner of a file can revoke access", func() {
			//1. Alice owns a file named "file.txt" and shares the file with Bob.
			alice.StoreFile("file.txt", []byte("content"))
			invitationPtr, _ := alice.CreateInvitation("file.txt", bobUsername)
			bob.AcceptInvitation(aliceUsername, invitationPtr, "from_alice.txt")

			//2. Alice also shares the file with Nilufar, and nilufar accepts invite
			invitationPtr2, _ := alice.CreateInvitation("file.txt", nilufarUsername)
			nilufar.AcceptInvitation(aliceUsername, invitationPtr2, "from_alice.txt")

			//3. Next, Bob tried to revoke Alice’s permission to this file:
			e1 := bob.RevokeAccess("from_alice.txt", aliceUsername)
			Expect(e1).ToNot(BeNil(), "Ownly owner of file should be able to revoke access.")

			//4. Nilufar tries to revoke Bob’s permission to this file:
			e2 := nilufar.RevokeAccess("from_alice.txt", bobUsername)
			Expect(e2).ToNot(BeNil(), "Ownly owner of file should be able to revoke access.")
		})

		//The client MUST allow the access to be revoked from another user even if the other user
		//does not have an active user session created by GetUser() (i.e. the user is currently offline
		It("should be able to revoke acces to a user when user is offline", func() {

		})

		//When owner revokes access to a user, ALL USERS WITHIN THAT SHARING BRANCH should have revoked access
		It("should be able to revoke acces to an entire branch of users", func() {

			//Alice -> Bob -> Olga
			// Bob -> Marco
			//Alice -> Nullifer
			alice.StoreFile(someFilename, someShortFileContent)

			/************************ CREATING TREE-LIKE STRUCTURE *********************************/
			//A-> B
			shareFileInfoPtr, err := alice.CreateInvitation(someFilename, bobUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Bob.")
			err1 := bob.AcceptInvitation(aliceUsername, shareFileInfoPtr, someOtherFilename)
			Expect(err1).To(BeNil(), "Bob could not accept Alice's shared invitation.")

			//A -> N
			shareFileInfoPtr2, err := alice.CreateInvitation(someFilename, nilufarUsername)
			Expect(err).To(BeNil(), "Alice failed to share a file with Nilufar.")
			err2 := nilufar.AcceptInvitation(aliceUsername, shareFileInfoPtr2, someOtherFilename)
			Expect(err2).To(BeNil(), "Nilufar could not accept Alice's shared invitation.")

			//B -> O
			shareFileInfoPtr3, err := bob.CreateInvitation(someOtherFilename, olgaUsername)
			Expect(err).To(BeNil(), "Bob failed to share a file with Olga.")
			err3 := olga.AcceptInvitation(bobUsername, shareFileInfoPtr3, someOtherFilename)
			Expect(err3).To(BeNil(), "Olga could not accept Bob's shared invitation.")

			//B -> Marco
			shareFileInfoPtr4, err := bob.CreateInvitation(someOtherFilename, marcoUsername)
			Expect(err).To(BeNil(), "Bob failed to share a file with Marco.")
			err4 := marco.AcceptInvitation(bobUsername, shareFileInfoPtr4, someOtherFilename)
			Expect(err4).To(BeNil(), "Marco could not accept Bob's shared invitation.")

			/************************ CUTTING OFF A SHARING BRANCH *********************************/

			//Alice revokes usage to Bob -> Bob, Olga, and Marco now should not be able to access & edit file
			e1 := alice.RevokeAccess(someFilename, bobUsername)
			Expect(e1).To(BeNil(), "Alice did not revoke access to Bob correctly.")

			//Bob, Marco, and Olga cannot acess & edit file
			//Testing Revoke Access for Bob
			_, e5 := bob.LoadFile(someOtherFilename)
			Expect(e5).ToNot(BeNil(), "Revoked User should NOT be able to Load to file.")

			e6 := bob.StoreFile(someOtherFilename, []byte("Testing, uno, dos, tres."))
			Expect(e6).ToNot(BeNil(), "Revoked User should NOT be able to Store to file.")

			e7 := bob.AppendToFile(someOtherFilename, []byte("some more stuff"))
			Expect(e7).ToNot(BeNil(), "Revoked User should NOT be able to Append to file.")

			//Testing Revoke Access for marco
			_, e8 := marco.LoadFile(someOtherFilename)
			Expect(e8).ToNot(BeNil(), "Revoked User should NOT be able to Load to file.")

			e9 := marco.StoreFile(someOtherFilename, []byte("Testing, uno, dos, tres."))
			Expect(e9).ToNot(BeNil(), "Revoked User should NOT be able to Store to file.")

			e10 := marco.AppendToFile(someOtherFilename, []byte("some more stuff"))
			Expect(e10).ToNot(BeNil(), "Revoked User should NOT be able to Append to file.")

			//Testing Revoke Access for Olga
			_, e11 := olga.LoadFile(someOtherFilename)
			Expect(e11).ToNot(BeNil(), "Revoked User should NOT be able to Load to file.")

			e12 := olga.StoreFile(someOtherFilename, []byte("Testing, uno, dos, tres."))
			Expect(e12).ToNot(BeNil(), "Revoked User should NOT be able to Store to file.")

			e13 := olga.AppendToFile(someOtherFilename, []byte("some more stuff"))
			Expect(e13).ToNot(BeNil(), "Revoked User should NOT be able to Append to file.")

			//Testing User (i.e. Nilufar) who should still have access
			//Overwritten + Loading Test
			e14 := nilufar.StoreFile(someOtherFilename, []byte("new stuff"))
			Expect(e14).To(BeNil(), "Owner should still be able to Store file.")

			nilufar_data, e2 := nilufar.LoadFile(someOtherFilename)
			Expect(e2).To(BeNil(), "Owner should still be able to Load file.")
			Expect("new stuff").To(BeEquivalentTo(string(nilufar_data)),
				"File should be overwritten after StoreFile() is called.")

			//Append test
			e15 := nilufar.AppendToFile(someOtherFilename, []byte("some more stuff"))
			Expect(e15).To(BeNil(), "Owner should still be able to Append to file.")

			nilufar_data2, _ := nilufar.LoadFile(someOtherFilename)
			real_data := "new stuff" + "some more stuff"
			Expect(real_data).To(BeEquivalentTo(string(nilufar_data2)),
				"File should be the same.")
		})

		//Test for revoked user trying to access file using Datastore directly
		It("should ensure integrity of file after a user is revoked", func() {
			//1. Alice owns a file named "file.txt" and shares the file with Bob.
			alice.StoreFile("file.txt", []byte("content"))
			invitationPtr, _ := alice.CreateInvitation("file.txt", bobUsername)
			bob.AcceptInvitation(aliceUsername, invitationPtr, "from_alice.txt")

			//2. Alice revokes access from Bob
			e1 := alice.RevokeAccess("file.txt", bobUsername)
			Expect(e1).To(BeNil(), "Ownly owner of file should be able to revoke access.")

			//3. Bob becomes a malicious user and tries to access file using Datastore and overwrite it/mess with its integrity
			datastore := userlib.DatastoreGetMap()
			datastore[invitationPtr] = []byte("gotcha")

			//4. Alice tries to access file -> File should be safe from malicious attack against file's integrity
			file_contents, e2 := alice.LoadFile("file.txt")
			Expect(e2).To(BeNil(), "File integrity should still be intact after malicious attack.")
			Expect([]byte("content")).To(BeEquivalentTo(string(file_contents)),
				"File contents should be the same.")

		})
	})
})
