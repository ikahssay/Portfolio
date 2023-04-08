package client

// CS 161 Project 2

// You MUST NOT change these default imports. ANY additional imports
// may break the autograder!

import (
	userlib "github.com/cs161-staff/project2-userlib"

	// Likewise, useful for debugging, etc.
	"encoding/hex"

	// Useful for string mainpulation.
	"strings"

	// Useful for formatting strings (e.g. `fmt.Sprintf`).
	"fmt"

	// Want to import errors.
	"errors"

	// Optional. You can remove the "_" there, but please do not touch
	// anything else within the import bracket.
	_ "strconv"

	/*TA,Shomil Jain, said we can add this import on Piazza!
	->" You can import bytes and use bytes.Equal or bytes.Compare (we added it to the list of allowed imports to the autograder,
		even though the starter code doesn’t import it).
	*/
	"bytes"
)

// This serves two purposes:
// a) It shows you some useful primitives, and
// b) it suppresses warnings for items not being imported.

// This function can be safely deleted!
func someUsefulThings() {
	// Creates a random UUID
	f := userlib.UUIDNew()
	userlib.DebugMsg("UUID as string:%v", f.String())

	// Example of writing over a byte of f
	f[0] = 10
	userlib.DebugMsg("UUID as string:%v", f.String())

	// Takes a sequence of bytes and renders as hex
	h := hex.EncodeToString([]byte("fubar"))
	userlib.DebugMsg("The hex: %v", h)

	// Marshals data into a JSON representation
	// Works well with Go structures!
	d, _ := userlib.Marshal(f)
	userlib.DebugMsg("The json data: %v", string(d))
	var g userlib.UUID
	userlib.Unmarshal(d, &g)
	userlib.DebugMsg("Unmashaled data %v", g.String())

	// errors.New(...) creates an error type!
	userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("This is an error")))

	// Generate a random private/public keypair.
	// The "_" indicates that we don't check for the error case here.
	var pk userlib.PKEEncKey
	var sk userlib.PKEDecKey
	pk, sk, _ = userlib.PKEKeyGen()
	userlib.DebugMsg("Key is %v, %v", pk, sk)

	// Useful for string interpolation.
	_ = fmt.Sprintf("%s_%d", "file", 1)
}

// User is the structure definition for a user record.
//Struct User {
//Username = “string”
//Password = [salt, password-encrypted(“salt || password”)]
//Dictionary: { hashed(“fileName1”) : [ownerName, UUID (points to the invitation struct, that then points to file info/struct)]
//}
type User struct {
	// You can add other fields here if you want...
	// Note for JSON to marshal/unmarshal, the fields need to
	// be public (start with a capital letter)

	Username string
	//MUST BE HASHED WITH SALT! -> MAKE SALT THE USERNAME CUZ USERNAME IS ALWAYS UNIQUE
	Password []byte
	//String of fileName must be hashed also!!
	Dict_of_Files map[string]Vertex
	Dict_of_keys  map[string]userlib.UUID
}

type Vertex struct {
	OwnerName                 string
	Uuid_to_invitation_struct userlib.UUID //Points to INVITATION STRUCT (that will point to the file)

	//If user isn’t the owner of this file, then have the dictionary be null (so only owner can have
	//access to the “sharing tree”)
	//UserAccessDict: {Alice: [ Bob('s invite UUID), Nullifer('s invite UUID)],
	//Bob: [Olga('s invite UUID), Macy('s invite UUID)]}
	UserAccessDict map[string][]userlib.UUID //Key is usercaller, value is username of people caller shared file to
}

type Invitation struct {
	OwnerName       string
	Users_sent_here []string     //Will be encrypted with public key
	Uuid_to_file    userlib.UUID //Will be encrypted with public key
}

//Whats inside each individual file
//FIX EVERYTHING THAT HAS TO DO WITH THIS NOW!
type File struct {
	//Filecontents -> encrypted + HMAC
	/* Encrypt-then-MAC CONTAINS TWO FIELDS: <ENC_k1(M), MAC_k2(NC_k1(M))>
	1. Encrypted file contents -> ENC_k1(M) -> So we can ACTUALLY SEE what the message/content is
	2/ MAC on the enc. file contents -> MAC_k2(NC_k1(M) -> So we can VERIFY message HASNT BEEN TAMPERED WITH
	*/
	Original_file_content_struct_uuid userlib.UUID
	Number_of_appends                 int
}

type Encrypted_Struct_or_Contents struct {
	Ciphertext []byte
	HMAC_tag   []byte
}

type Publicly_Encrypted_Contents struct {
	Ciphertext []byte
	DS_tag     []byte
}

func InitUser(username string, password string) (userdataptr *User, err error) {
	//TODO: This is a toy implementation.
	//var userdata User
	//userdata.Username = username

	//1. Usernames CANNOT be empty
	if len(username) == 0 {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Username cannot be empty")))
		return nil, errors.New(strings.ToTitle("Username cannot be empty"))
	}

	//2. Check if username already exists (need username to be unique) in DataStore
	username_in_bytes := []byte(username)
	Uuid_of_user, e := userlib.UUIDFromBytes(userlib.Hash(username_in_bytes))
	//Returns an error if the slice is less than length 16.
	if e != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("The byte slice needed to create UUID is less than length 16.")))
	}

	//If you found a userstruct with that name already in the Datastore, then username is not unique!
	//-> Error and tell user to choose another usernmae
	if _, found := userlib.DatastoreGet(Uuid_of_user); found {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("This username is already in our system. Please choose another.")))
		return nil, errors.New(strings.ToTitle("This username is already in our system. Please choose another."))
	} else {
		//4. Create a new User struct, initialize variables, and save it to the DataStore
		var userdata User
		userdata.Username = username

		//Argon2Key(password []byte, salt []byte, keyLen uint32) ([]byte)
		username_in_bytes := []byte(username)
		password_in_bytes := []byte(password)
		userdata.Password = userlib.Argon2Key(password_in_bytes, username_in_bytes, 48) //The usernmame is the salt (will always be unique), and password length is 48
		userdata.Dict_of_Files = make(map[string]Vertex)
		userdata.Dict_of_keys = make(map[string]userlib.UUID)

		//Create the user's public (and private) PKEnc and DS keys, and store that in Dict_of_Keys
		_, _, e1 := CreateAndSaveKeys(&userdata)
		if e1 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("CreateAndSaveKeys is NOT working correctly!")))
			return nil, errors.New(strings.ToTitle("CreateAndSaveKeys is NOT working correctly!"))
		}

		//Saving to datastore
		userdata_in_bytes, e2 := userlib.Marshal(userdata)
		userlib.DatastoreSet(Uuid_of_user, userdata_in_bytes)
		if e2 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Marshalled incorrectly.")))
		}

		return &userdata, nil
	}
}

func GetUser(username string, password string) (userdataptr *User, err error) {
	var userdata User

	//1. Username CANNOT be empty
	if len(username) == 0 {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Username cannot be empty.")))
		return nil, errors.New(strings.ToTitle("Username cannot be empty."))
	}

	//2. Check that username the already exists in DataStore, and
	//3. Check the passwrod, and get the User struct from DataStore and return it

	username_in_bytes := []byte(username)
	Uuid_of_user, e := userlib.UUIDFromBytes(userlib.Hash(username_in_bytes))
	//Returns an error if the slice is less than length 16.
	if e != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("The byte slice needed to create UUID is less than length 16.")))
	}

	//If user/username is not found in database, throw an error, if not
	if _, found := userlib.DatastoreGet(Uuid_of_user); !found {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("User does not exist with that username.")))
		return nil, errors.New(strings.ToTitle("User does not exist with that username."))
	} else {
		var user_struct_bytes, _ = userlib.DatastoreGet(Uuid_of_user)
		err2 := userlib.Unmarshal(user_struct_bytes, &userdata)
		if err2 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Unmarshalled incorrectly.")))
		}

		//Test to see if user wrote in correct password (password matches saved & encrypted password)
		username_in_bytes := []byte(username)
		password_in_bytes := []byte(password)
		real_password := userdata.Password
		testing_password := userlib.Argon2Key(password_in_bytes, username_in_bytes, 48)

		if are_equal := bytes.Equal(real_password, testing_password); !are_equal {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Password is incorrect.")))
			return nil, errors.New(strings.ToTitle("Password is incorrect."))
		} else {
			return &userdata, nil
		}
	}
}

/*	-Given a filename in the personal namespace of the caller/user, this function persistently stores the given content (stores
		it in the Datastore) for future retrieval using the same filename.

	- If the given filename already exists in the personal namespace of the caller/user, then the content of the corresponding
		file is OVERWRITTEN.
*/
func (userdata *User) StoreFile(filename string, content []byte) (err error) {

	/*TODO: This is a toy implementation.
	storageKey, _ := userlib.UUIDFromBytes([]byte(filename + userdata.Username))
	jsonData, _ := userlib.Marshal(content)
	userlib.DatastoreSet(storageKey, jsonData)
	*/

	if userdata == nil {
		return errors.New(strings.ToTitle("User has not logged in or was not created"))
	}

	//IN THE CASE where user is operating on the same file WITH 2 DIFFERENT USERS/DEVICES, SYNC/UPDATE USER STRUCT/POINTER
	username_in_bytes := []byte(userdata.Username)
	user_uuid, _ := userlib.UUIDFromBytes(userlib.Hash(username_in_bytes))

	syncing_user, _ := userlib.DatastoreGet(user_uuid)
	userlib.Unmarshal(syncing_user, userdata)

	//1. If this is a fresh user, and they've never saved a file before, create a new Dict_of_Files dictionary, and instantiate it
	if userdata.Dict_of_Files == nil {
		CreateFileMap(userdata, filename, content)
	} else {
		//2. If Dict_of_Files is already created, and this is a NEW file being stored, CREATE A FILEPATH and add the new file with its contents
		filename_in_bytes := []byte(filename)
		hashed_file := userlib.Hash(filename_in_bytes)
		file_flow, exists := userdata.Dict_of_Files[hex.EncodeToString(hashed_file)]

		//If filepath doesn't exist, then create new file and add it
		if !exists {
			//Create NEW file path and file
			_, e1 := AddNewFile(userdata, filename, &file_flow, content)
			if e1 != nil {
				userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not create and add new file correctly")))
				return errors.New(strings.ToTitle("Did not create and add new file correctly"))
			}
		} else {

			//3. If Dict_of_Files is already created, and this is an EXISTING file being stored, FOLLOW THE FILEPATH and OVERWRITE file with contents
			//Follow existing file path and access/update file
			invite_UUID := file_flow.Uuid_to_invitation_struct

			//Need to decrypt the invitation struct to access it and its file UUID/pointer

			//Getting encrypted struct (of invitation struct)
			encrypted_invite_struct_in_bytes, encrypted_invite_struct_exists := userlib.DatastoreGet(invite_UUID)
			var encrypted_struct Encrypted_Struct_or_Contents
			e9 := userlib.Unmarshal(encrypted_invite_struct_in_bytes, &encrypted_struct)
			if e9 != nil {
				userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Unmarshalled incorrectly.")))
				return errors.New(strings.ToTitle("Unmarshalled incorrectly.."))
			}

			//Decrepyt
			if !encrypted_invite_struct_exists {
				userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Do not have invitation struct/acess to file-> User is revoked")))
				return errors.New(strings.ToTitle("Do not have invitation struct/acess to file-> User is revoked"))
			} else {
				//Decrypt invitation struct
				var invite, e10 = DecryptStruct(userdata, &encrypted_struct, invite_UUID, filename)
				if e10 != nil {
					userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Didn't decrypt invitation struct correctly.")))
					return errors.New(strings.ToTitle("Didn't decrypt invitation struct correctly."))
				}

				//Go to the file struct (using the invitation struct)
				var file_struct File
				file_UUID := invite.Uuid_to_file
				file_struct_in_bytes, exists := userlib.DatastoreGet(file_UUID)
				e3 := userlib.Unmarshal(file_struct_in_bytes, &file_struct)
				if !exists {
					userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Stored file incorrectly in DataStore.")))
					return errors.New(strings.ToTitle("Stored file incorrectly in DataStore."))
				}
				if e3 != nil {
					userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Unmarshalled incorrectly.")))
					return errors.New(strings.ToTitle("Unmarshalled incorrectly."))
				}

				//ii. Overwrite file contents in file
				new_file_contents_UUID := userlib.UUIDNew()
				file_struct.Original_file_content_struct_uuid = new_file_contents_UUID
				file_struct.Number_of_appends = 0

				encrypted_content, encrpyt_and_HMAC_content, e6 := EncryptContents(userdata, invite, content)
				encrypted_file_content_struct := &Encrypted_Struct_or_Contents{
					encrypted_content,
					encrpyt_and_HMAC_content,
				}

				if e6 != nil {
					userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not encrypt-then_HMAC file contents correctly!")))
					return errors.New(strings.ToTitle("Did not encrypt-then_HMAC file contents correctly!"))
				}

				//4. SAVE EVERYTHING!
				//A.  Save User struct (and any changes that came with it - like adding a new file or dict-of-files)
				//Getting UUID of user
				username_in_bytes := []byte(userdata.Username)
				Uuid_of_user, e := userlib.UUIDFromBytes(userlib.Hash(username_in_bytes))
				if e != nil {
					userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("The byte slice needed to create UUID is less than length 16.")))
				}

				//Saving to datastore
				userdata_in_bytes, e2 := userlib.Marshal(userdata)
				userlib.DatastoreSet(Uuid_of_user, userdata_in_bytes)
				if e2 != nil {
					userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Marshalled incorrectly.")))
				}

				//B. SaveInvitation struct

				//i.Encrypt + HMAC entire Invitation struct so only non-malicious users can see it
				//Creating key for encrypting invitationStruct (once it's created)
				encrypted_struct, encrpyt_and_HMAC_inv_struct, err := EncryptInvitationStruct(userdata, invite, invite_UUID)
				if err != nil {
					userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not encrypt-then_HMAC invitation struct correctly!")))
					return errors.New(strings.ToTitle("Did not encrypt-then_HMAC correctly!"))
				}

				//Auh_enc_inv_struct = { ENC_k1(inv_struct_in_bytes), MAC_k2(ENC_k1(inv_struct_in_bytes))}
				auth_enc_struct := &Encrypted_Struct_or_Contents{
					encrypted_struct,
					encrpyt_and_HMAC_inv_struct,
				}

				//ii. Save to DataStore
				auth_enc_inv_struct_bytes, e := userlib.Marshal(auth_enc_struct)
				if e != nil {
					return errors.New(strings.ToTitle("Paniiiccccc"))
				}

				userlib.DatastoreSet(invite_UUID, auth_enc_inv_struct_bytes)

				//C. Save updated file struct (and its file contents)
				//i. Save File in DataStore
				// -> Need fileUUID & its corresponding file struct

				file_struct_in_bytes, err2 := userlib.Marshal(file_struct)
				userlib.DatastoreSet(invite.Uuid_to_file, file_struct_in_bytes)
				if err2 != nil {
					userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not File struct correctly.")))
					return errors.New(strings.ToTitle("Did not File struct correctly."))
				}

				//D. Save new encrypted_file_contents_struct
				file_content_struct_in_bytes, err := userlib.Marshal(encrypted_file_content_struct)
				userlib.DatastoreSet(new_file_contents_UUID, file_content_struct_in_bytes)
				if err != nil {
					userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not save encrypted file contents struct correctly.")))
					return errors.New(strings.ToTitle("Did not save encrypted file contents struct correctly."))
				}

				return nil
			}
		}
	}

	return nil
}

//Helper Function: Encrypt + HMAC ALL STRUCTS before saving them onto DataStore! <- Is that necessary though?
func EncryptInvitationStruct(userdata *User, invite *Invitation, invite_UUID userlib.UUID) (encrypted_struct []byte, encrpyt_and_HMAC_inv_struct []byte, err error) {

	//Creating HMAC and encrypted keys using inv. struct's UUID
	inv_struct_UUID_in_bytes, e3 := userlib.Marshal(invite_UUID)
	incremented_UUID := userlib.Hash([]byte(hex.EncodeToString(inv_struct_UUID_in_bytes) + "100"))
	symmetric_enc_key := userlib.Argon2Key(inv_struct_UUID_in_bytes, incremented_UUID, 16)
	HMAC_key, err := userlib.HashKDF(inv_struct_UUID_in_bytes[:16], []byte("invitation struct's HMAC key"))

	if e3 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Marshalled incorrectly.")))
		return nil, nil, errors.New(strings.ToTitle("Marshalled incorrectly."))
	}
	if err != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Created invitation_struct_key incorrectly.")))
		return nil, nil, errors.New(strings.ToTitle("Created invitation_struct_key incorrectly."))
	}

	//Encrypting + HMAC Inv. Struct
	invite_to_bytes, e4 := userlib.Marshal(invite)
	iv := userlib.RandomBytes(16)
	if e4 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Marshalled incorrectly.")))
		return nil, nil, errors.New(strings.ToTitle("Marshalled incorrectly."))
	}

	encrypted_struct = userlib.SymEnc(symmetric_enc_key[:16], iv, invite_to_bytes)
	encrpyt_and_HMAC_struct, e5 := userlib.HMACEval(HMAC_key[:16], encrypted_struct)
	if e5 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not encrypt-then_HMAC correctly!")))
		return nil, nil, errors.New(strings.ToTitle("Did not encrypt-then_HMAC correctly!"))
	}

	return encrypted_struct, encrpyt_and_HMAC_struct, nil
}

//Helper Function
func DecryptStruct(userdata *User, encrypted_struct *Encrypted_Struct_or_Contents, invite_UUID userlib.UUID, filename string) (invite *Invitation, err error) {
	//Creating HMAC and encrypted keys using inv. struct's UUID
	inv_struct_UUID_in_bytes, e3 := userlib.Marshal(invite_UUID)
	incremented_UUID := userlib.Hash([]byte(hex.EncodeToString(inv_struct_UUID_in_bytes) + "100"))
	symmetric_enc_key := userlib.Argon2Key(inv_struct_UUID_in_bytes, incremented_UUID, 16)
	HMAC_key, err := userlib.HashKDF(inv_struct_UUID_in_bytes[:16], []byte("invitation struct's HMAC key"))

	if e3 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Marshalled incorrectly.")))
		return nil, errors.New(strings.ToTitle("Marshalled incorrectly."))
	}
	if err != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Created invitation_struct_key incorrectly.")))
		return nil, errors.New(strings.ToTitle("Created invitation_struct_key incorrectly."))
	}

	//Check HMAC_tag matches encrypted struct's HMAC_tag
	cipher_text := encrypted_struct.Ciphertext
	tester_tag := encrypted_struct.HMAC_tag
	real_tag, err := userlib.HMACEval(HMAC_key[:16], cipher_text)
	if err != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not encrypt-then_HMAC correctly!")))
		return nil, errors.New(strings.ToTitle("Did not encrypt-then_HMAC correctly!"))
	}
	//Checking if they're equal
	//HMACEqual(a []byte, b []byte) (bool)
	the_same := userlib.HMACEqual(tester_tag, real_tag)
	if !the_same {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("HMAC_tags don't match!!!")))
		return nil, errors.New(strings.ToTitle("HMAC_tags don't match!!!"))
	}

	//If it matches, then decrypt and return (true/decrypted) invitation struct
	inv_struct_in_bytes := userlib.SymDec(symmetric_enc_key, cipher_text)

	var inv_struct *Invitation
	err2 := userlib.Unmarshal(inv_struct_in_bytes, &inv_struct)
	if err2 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Unmarshalled incorrectly.")))
		return nil, errors.New(strings.ToTitle("Unmarshalled incorrectly."))
	}

	return inv_struct, nil
}

//Helper Function
func EncryptContents(userdata *User, invite *Invitation, content []byte) (encrypted_content []byte, encrpyt_and_HMAC_content []byte, err error) {
	//Creating HMAC and encrypted keys using file struct's UUID
	fileUUID_in_bytes, e5 := userlib.Marshal(invite.Uuid_to_file)
	incremented_fileUUID_in_bytes := userlib.Hash([]byte(hex.EncodeToString(fileUUID_in_bytes) + "100"))
	iv := userlib.RandomBytes(16)

	file_symmetric_key := userlib.Argon2Key(fileUUID_in_bytes, incremented_fileUUID_in_bytes, 16)
	file_HMAC_key, e6 := userlib.HashKDF(fileUUID_in_bytes[:16], []byte("file's HMAC key"))

	if e5 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Marshalled fileUUID incorrectly.")))
		return nil, nil, errors.New(strings.ToTitle("Marshalled incorrectly."))
	}
	if e6 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Created HKDF for file_MAc key incorectly.")))
		return nil, nil, errors.New(strings.ToTitle("Created HKDF inccorectly."))
	}

	//Peform encrpyt-then-MAC
	encrypted_content = userlib.SymEnc(file_symmetric_key, iv, content)
	encrpyt_and_HMAC_content, e7 := userlib.HMACEval(file_HMAC_key[:16], encrypted_content)
	if e7 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not encrypt-then_HMAC correctly!")))
		return nil, nil, errors.New(strings.ToTitle("Did not encrypt-then_HMAC correctly!"))
	}

	//Return <Enc_k1(file_contents), HMAC_k2(Enc_k1(file_contents))>
	return encrypted_content, encrpyt_and_HMAC_content, nil
}

//Helper Function: Decrypts file content. If there are APPENDED FILE CONTENTS, it will also decrypte ALL of them, and return
//all of the decrypted content in one package
func DecryptContents(userdata *User, invite *Invitation) (content []byte, err error) {
	//Get File struct
	var file_struct File
	file_struct_in_bytes, exist := userlib.DatastoreGet(invite.Uuid_to_file)
	e2 := userlib.Unmarshal(file_struct_in_bytes, &file_struct)

	if !exist {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Couldn't get file from DataStore.")))
		return nil, errors.New(strings.ToTitle("Couldn't get file from DataStore."))
	}
	if e2 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Unmarshalled file incorrectly.")))
		return nil, errors.New(strings.ToTitle("Unmarshalled file incorrectly."))
	}

	//Creating HMAC and encrypted keys using file struct's UUID
	fileUUID_in_bytes, e5 := userlib.Marshal(invite.Uuid_to_file)
	incrementing_fileUUID_in_bytes := userlib.Hash([]byte(hex.EncodeToString(fileUUID_in_bytes) + "100"))
	file_symmetric_key := userlib.Argon2Key(fileUUID_in_bytes, incrementing_fileUUID_in_bytes, 16)
	file_HMAC_key, e6 := userlib.HashKDF(fileUUID_in_bytes[:16], []byte("file's HMAC key"))

	if e5 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Marshalled fileUUID incorrectly.")))
		return nil, errors.New(strings.ToTitle("Marshalled incorrectly."))
	}
	if e6 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Created HKDF for file_MAc key incorectly.")))
		return nil, errors.New(strings.ToTitle("Created HKDF inccorectly."))
	}

	//Access each APPENDED file_content struct, check that the HMAC_tag matches encrypted struct's HMAC_tag, then decrypt each file content
	//Save decrypted contents in one []byte (APPEND). Put that in a new struct
	if file_struct.Number_of_appends > 0 {
		var original_contents []byte
		var concatentated_contents []byte

		//For each file_content struct,
		for i := 1; i < file_struct.Number_of_appends+2; i++ {
			//Retrieving file_content_struct of original content and appended content -> will decrypt original content LAST
			//What below is saying: incremented_uuid := UUIDFromBytes(Hash(old/original_uuid.String() + "1")))
			var file_content_struct_uuid userlib.UUID
			if i != file_struct.Number_of_appends+1 {
				temp := []byte(file_struct.Original_file_content_struct_uuid.String() + string(rune(i)))
				file_content_struct_uuid, err = userlib.UUIDFromBytes(userlib.Hash(temp))
				if err != nil {
					userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not deterministically create appended_file_contents' UUIDs!")))
					return nil, errors.New(strings.ToTitle("Did not deterministically create appended_file_contents' UUIDs!"))
				}
			} else {
				file_content_struct_uuid = file_struct.Original_file_content_struct_uuid
			}

			file_content_struct_in_bytes, exists := userlib.DatastoreGet(file_content_struct_uuid)
			if !exists {
				userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("File_content_struct was saved in DataStore incorectly.")))
				return nil, errors.New(strings.ToTitle("File_content_struct was saved in DataStore incorectly."))
			}

			var file_content_struct Encrypted_Struct_or_Contents
			err2 := userlib.Unmarshal(file_content_struct_in_bytes, &file_content_struct)
			if err2 != nil {
				userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("File_content_struct was unmarshalled incorectly.")))
				return nil, errors.New(strings.ToTitle("File_content_struct was unmarshalled incorectly."))
			}

			//Checking that the HMAC_tag matches encrypted struct's HMAC_tag
			cipher_text := file_content_struct.Ciphertext
			tester_tag := file_content_struct.HMAC_tag
			real_tag, err3 := userlib.HMACEval(file_HMAC_key[:16], cipher_text)
			the_same := userlib.HMACEqual(tester_tag, real_tag)

			if !the_same {
				userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("HMAC_tags don't match!!!")))
				return nil, errors.New(strings.ToTitle("HMAC_tags don't match!!!"))
			}
			if err3 != nil {
				userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not encrypt-then_HMAC correctly!")))
				return nil, errors.New(strings.ToTitle("Did not encrypt-then_HMAC correctly!"))
			}

			//If it matches, then decrypt file contents and append it to concatentated_contents
			//Instead of decrypting the original contents (and its like 1TB) and then appending each piece of new file contents to it,
			//Decrypt the original contents ONCE, and save it. Then, append all new file_contents and save it somewhere else.
			//At the end, append everything.
			file_contents := userlib.SymDec(file_symmetric_key, cipher_text)

			/**	How to append two byte slices, where a is []byte and b is []byte
				-> Instead of this: for i := 0; i < len(b); i++ {
										a = append(a,b[i])
									}
			    -> You can do a quicker way! (And it does the same thing as above)
					a= append(a, b...)
			*/
			if i == file_struct.Number_of_appends+1 {
				original_contents = file_contents
			} else {
				concatentated_contents = append(concatentated_contents, file_contents...)
			}

		}

		//Get original contents and append concatenated_contents to it
		appended_file_contents := append(original_contents, concatentated_contents...)

		//Send back decrypted and appended file contents
		return appended_file_contents, nil

	} else {
		//Retrieving the ONLY ONE file_content_struct
		file_content_struct_in_bytes, exists := userlib.DatastoreGet(file_struct.Original_file_content_struct_uuid)
		if !exists {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("File_content_struct was saved in DataStore incorectly.")))
			return nil, errors.New(strings.ToTitle("File_content_struct was saved in DataStore incorectly."))
		}

		var file_content_struct Encrypted_Struct_or_Contents
		err2 := userlib.Unmarshal(file_content_struct_in_bytes, &file_content_struct)
		if err2 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("File_content_struct was unmarshalled incorectly.")))
			return nil, errors.New(strings.ToTitle("File_content_struct was unmarshalled incorectly."))
		}

		//Checking that the HMAC_tag matches encrypted struct's HMAC_tag
		cipher_text := file_content_struct.Ciphertext
		tester_tag := file_content_struct.HMAC_tag
		real_tag, err3 := userlib.HMACEval(file_HMAC_key[:16], cipher_text)
		the_same := userlib.HMACEqual(tester_tag, real_tag)

		if !the_same {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("HMAC_tags don't match!!!")))
			return nil, errors.New(strings.ToTitle("HMAC_tags don't match!!!"))
		}
		if err3 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not encrypt-then_HMAC correctly!")))
			return nil, errors.New(strings.ToTitle("Did not encrypt-then_HMAC correctly!"))
		}

		//If it matches, then decrypt file contents and return it
		file_contents := userlib.SymDec(file_symmetric_key, cipher_text)
		return file_contents, nil
	}
}

//Helper Function: Create File map
func CreateFileMap(userdata *User, filename string, content []byte) {
	//1. Create Dict_of_Files map[string] Vertex{}, and initialize it
	new_dict := make(map[string]Vertex)
	userdata.Dict_of_Files = new_dict
	userdata.Dict_of_keys = make(map[string]userlib.UUID)

	//2. Create File struct where File contents will be stored, and store the contents
	filename_in_bytes := []byte(filename)
	hashed_file := userlib.Hash(filename_in_bytes)
	file_flow := userdata.Dict_of_Files[hex.EncodeToString(hashed_file)]
	AddNewFile(userdata, filename, &file_flow, content)
}

//Helper Function: Create NEW file path and add a NEW file to Dict_of_files.
//Then, save encrypted + MAC contents in new file.
//Then, save everything in DataStore
func AddNewFile(userdata *User, filename string, file_flow *Vertex, content []byte) (file_struct *File, err error) {
	/* What I'm trying to create:
	Dictionary of Files: { hashed(“fileName1”): [ Ownername of file (string),
													UUID to the invitation struct (that will point us to the file),
													UserAccessDict {Alice: Bob, Nullifer,
																	Bob: Olga, Macy} ]
	*/

	//1. Create FilePath
	filename_in_bytes := []byte(filename)
	hashed_file := userlib.Hash(filename_in_bytes)

	//Things to put in File's Vertex
	invitation_struct_UUID := userlib.UUIDNew()
	user_access_dict := make(map[string][]userlib.UUID)
	user_access_dict[userdata.Username] = make([]userlib.UUID, 0) //Says: {Alice: []}

	file_flow = &Vertex{
		userdata.Username,
		invitation_struct_UUID,
		user_access_dict, //Rn, user_access_dict looks like: {Alice: ""}
	}

	//Add filepath in Dict_of_Files
	userdata.Dict_of_Files[hex.EncodeToString(hashed_file)] = *file_flow

	//2. Create Invitation struct (that leads to file) and initialize it
	//Need the key to decrypt info -> SYMMETRIC KEY-> HKDF{UUID to invitation struct}
	invite, e1 := CreateInvitationStruct(userdata, filename, file_flow)
	if e1 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Created invitation struct incorrectly.")))
		return nil, errors.New(strings.ToTitle("Created invitation struct incorrectly."))
	}

	//3. Create File struct where File contents will be stored, and store the contents
	/* Filecontent must be Encrypt-then-HMAC! Keys are created using the file's UUID
			-Use Password-Based Key Derivation Function (PKDF) (with the salt being the file owner’s name) to create
				the symmetric key encryption
	`			-Use Hash-Based Key Derivation Function (HKDF) to create the key for HMAC
	*/
	//Encrypt contents
	encrypted_content, encrpyt_and_HMAC_content, e := EncryptContents(userdata, invite, content)
	if e != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Marshalled incorrectly.")))
		return nil, errors.New(strings.ToTitle("Marshalled incorrectly."))
	}

	//ii. Save content into a new file_content_struct. Then save the file_content_struct into (new) File's list of UUIDs
	new_file_content_struct := &Encrypted_Struct_or_Contents{
		encrypted_content,
		encrpyt_and_HMAC_content,
	}

	new_file_content_UUID := userlib.UUIDNew()
	new_file := &File{
		new_file_content_UUID,
		0,
	}

	//4. Save everything
	//i. Save FilePath and (and added File) by updating userdata in DataStore
	username_in_bytes := []byte(userdata.Username)
	Uuid_of_user, e := userlib.UUIDFromBytes(userlib.Hash(username_in_bytes))

	if e != nil { //Returns an error if the slice is less than length 16.
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("The byte slice needed to create UUID is less than length 16.")))
		return nil, errors.New(strings.ToTitle("The byte slice needed to create UUID is less than length 16."))
	}

	userdata_in_bytes, e2 := userlib.Marshal(userdata)
	userlib.DatastoreSet(Uuid_of_user, userdata_in_bytes)
	if e2 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Marshalled incorrectly.")))
		return nil, errors.New(strings.ToTitle("Marshalled incorrectly."))
	}

	//ii. Save new file_content_struct in DataStore
	new_file_content_struct_in_bytes, e4 := userlib.Marshal(new_file_content_struct)
	userlib.DatastoreSet(new_file_content_UUID, new_file_content_struct_in_bytes)
	if e4 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Marshalled incorrectly.")))
		return nil, errors.New(strings.ToTitle("Marshalled incorrectly."))
	}

	//iii. Save new File (and its contents/UUIDs) in DataStore
	// -> Need fileUUID & its corresponding file struct
	file_struct_in_bytes, e3 := userlib.Marshal(new_file)
	userlib.DatastoreSet(invite.Uuid_to_file, file_struct_in_bytes)
	if e3 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Marshalled incorrectly.")))
		return nil, errors.New(strings.ToTitle("Marshalled incorrectly."))
	}

	return new_file, nil
}

//Helper Function
//input: FilePath/Vertex pointer
//output: invitation_struct_pointer
func CreateInvitationStruct(userdata *User, filename string, file_flow *Vertex) (invite *Invitation, err error) {

	//Creating invitation struct
	invitation_struct_UUID := file_flow.Uuid_to_invitation_struct
	filename_in_bytes := []byte(filename)
	hashed_filename := userlib.Hash(filename_in_bytes)
	fileUUID := userlib.UUIDNew()
	users_sent_here := make([]string, 0) //Creating an empty array of strings [""]
	users_sent_here = append(users_sent_here, userdata.Username)

	invite = &Invitation{ //Will be encrypted with symmetric key = HKDF{UUID to invitation struct}
		userdata.Dict_of_Files[hex.EncodeToString(hashed_filename)].OwnerName,
		users_sent_here, //if owner is creating new file, should only be ["Alice"]
		fileUUID,
	}

	//Encrypt entire Invitation struct so only non-malicious user can see it
	encrypted_struct, encrpyt_and_HMAC_inv_struct, err := EncryptInvitationStruct(userdata, invite, invitation_struct_UUID)
	if err != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not encrypt-then_HMAC invitation struct correctly!")))
		return nil, errors.New(strings.ToTitle("Did not encrypt-then_HMAC correctly!"))
	}

	//Auh_enc_inv_struct = { ENC_k1(inv_struct_in_bytes), MAC_k2(ENC_k1(inv_struct_in_bytes))}
	auth_enc_struct := &Encrypted_Struct_or_Contents{
		encrypted_struct,
		encrpyt_and_HMAC_inv_struct,
	}

	//ii. Save to DataStore
	auth_enc_inv_struct_bytes, e := userlib.Marshal(auth_enc_struct)
	userlib.DatastoreSet(invitation_struct_UUID, auth_enc_inv_struct_bytes)
	if e != nil {
		return nil, errors.New(strings.ToTitle("Paniiiccccc"))
	}

	return invite, nil
}

func (userdata *User) AppendToFile(filename string, new_content []byte) error {
	if userdata == nil {
		return errors.New(strings.ToTitle("User has not logged in or was not created"))
	}

	//IN THE CASE where user is operating on the same file WITH 2 DIFFERENT USERS/DEVICES, SYNC/UPDATE USER STRUCT/POINTER
	username_in_bytes := []byte(userdata.Username)
	user_uuid, _ := userlib.UUIDFromBytes(userlib.Hash(username_in_bytes))

	syncing_user, _ := userlib.DatastoreGet(user_uuid)
	userlib.Unmarshal(syncing_user, userdata)

	//1. Check if hashed fileName exists in Dict_of_Files
	filename_in_bytes := []byte(filename)
	hashed_file := userlib.Hash(filename_in_bytes)
	file_flow, exists := userdata.Dict_of_Files[hex.EncodeToString(hashed_file)]

	//2. Follow path to file. If path to file is missing, throw an error
	if !exists { //If filepath doesn't exist, then throw an error
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("File with that file name does not exist")))
		return errors.New(strings.ToTitle("File with that file name does not exist"))
	} else {
		//Following path to file
		//Follow existing file path and access/update file
		invite_UUID := file_flow.Uuid_to_invitation_struct

		//Need to decrypt the invitation struct to access it and its file UUID/pointer

		//Getting encrypted struct (of invitation struct)
		var encrypted_struct Encrypted_Struct_or_Contents
		encrypted_invite_struct_in_bytes, encrypted_invite_struct_exists := userlib.DatastoreGet(invite_UUID)
		e9 := userlib.Unmarshal(encrypted_invite_struct_in_bytes, &encrypted_struct)
		if e9 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Unmarshalled incorrectly.")))
			return errors.New(strings.ToTitle("Unmarshalled incorrectly.."))
		}

		//Decrypt
		if !encrypted_invite_struct_exists { //If encrypted_struct doesn't exist in DataStore
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Do not have invitation struct/acess to file-> User is revoked")))
			return errors.New(strings.ToTitle("Do not have invitation struct/acess to file-> User is revoked"))
		}
		//Decrypt invitation struct
		var invite, e10 = DecryptStruct(userdata, &encrypted_struct, invite_UUID, filename)
		if e10 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Decrypted invitation struct incorrectly.")))
			return errors.New(strings.ToTitle("Decrypted invitation structincorrectly."))
		}

		// 3. Get the existing encrypted file contents, encrypt new file contents, put it in a new file_contents srtuct, save that new
		//file_content_struct's UUID in existing File's list of UUIDs.
		//i. Encrypt new file content and put it to new file_content struct,
		//Encrypting contents
		new_encrypted_content, new_encrpyt_and_HMAC_content, e := EncryptContents(userdata, invite, new_content)
		if e != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Marshalled incorrectly.")))
			return errors.New(strings.ToTitle("Marshalled incorrectly."))
		}

		//Saving content into a new file_content_struct.
		new_file_content_struct := &Encrypted_Struct_or_Contents{
			new_encrypted_content,
			new_encrpyt_and_HMAC_content,
		}

		//ii. Add that new struct's UUID to (existing) File's list of file_content_UUIDs

		//Getting File struct
		var file_struct File
		file_struct_in_bytes, exist := userlib.DatastoreGet(invite.Uuid_to_file)
		e2 := userlib.Unmarshal(file_struct_in_bytes, &file_struct)
		if !exist {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Couldn't get file from DataStore.")))
			return errors.New(strings.ToTitle("Couldn't get file from DataStore."))
		}
		if e2 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Unmarshalled file incorrectly.")))
			return errors.New(strings.ToTitle("Unmarshalled file incorrectly."))
		}

		//Deterministically creating the new UUID for the new/appended file_content_struct
		file_struct.Number_of_appends = file_struct.Number_of_appends + 1
		temp := []byte(file_struct.Original_file_content_struct_uuid.String() + string(rune(file_struct.Number_of_appends)))
		new_file_content_UUID, err := userlib.UUIDFromBytes(userlib.Hash(temp))
		if err != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not deterministically create appended_file_contents' UUIDs!")))
			return errors.New(strings.ToTitle("Did not deterministically create appended_file_contents' UUIDs!"))
		}

		//4. Save/update File in Datastore and new file_contents struct in DataStore
		//i. Save updated File
		file_struct_in_bytes, e11 := userlib.Marshal(file_struct)
		userlib.DatastoreSet(invite.Uuid_to_file, file_struct_in_bytes)
		if e11 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Marshalled incorrectly.")))
			return errors.New(strings.ToTitle("Marshalled incorrectly."))
		}

		//ii. Save new file_content_struct
		new_file_content_struct_in_bytes, e12 := userlib.Marshal(new_file_content_struct)
		userlib.DatastoreSet(new_file_content_UUID, new_file_content_struct_in_bytes)
		if e12 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Marshalled incorrectly.")))
			return errors.New(strings.ToTitle("Marshalled incorrectly."))
		}

		return nil
	}
}

func (userdata *User) LoadFile(filename string) (content []byte, err error) {

	if userdata == nil {
		return nil, errors.New(strings.ToTitle("User has not logged in or was not created"))
	}

	//IN THE CASE where user is operating on the same file WITH 2 DIFFERENT USERS/DEVICES, SYNC/UPDATE USER STRUCT/POINTER
	username_in_bytes := []byte(userdata.Username)
	user_uuid, _ := userlib.UUIDFromBytes(userlib.Hash(username_in_bytes))

	syncing_user, _ := userlib.DatastoreGet(user_uuid)
	userlib.Unmarshal(syncing_user, userdata)

	//1. Check if hashed fileName exists in Dict_of_Files
	filename_in_bytes := []byte(filename)
	hashed_file := userlib.Hash(filename_in_bytes)
	var file_flow, exists = userdata.Dict_of_Files[hex.EncodeToString(hashed_file)]

	//2. Follow path to file. If path to file is missing, throw an error
	//If filepath doesn't exist, then throw an error
	if !exists {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("File with that file name does not exist")))
		return nil, errors.New(strings.ToTitle("File with that file name does not exist"))
	} else {
		//Following path to file
		//Follow existing file path and access/update file
		invite_UUID := file_flow.Uuid_to_invitation_struct

		//Need to decrypt the invitation struct to access it and its file UUID/pointer

		//Getting encrypted struct (of invitation struct)
		encrypted_invite_struct_in_bytes, encrypted_invite_struct_exists := userlib.DatastoreGet(invite_UUID)
		var encrypted_struct Encrypted_Struct_or_Contents
		e9 := userlib.Unmarshal(encrypted_invite_struct_in_bytes, &encrypted_struct)
		if e9 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Unmarshalled incorrectly.")))
			return nil, errors.New(strings.ToTitle("Unmarshalled incorrectly.."))
		}

		//Decrepyt
		if !encrypted_invite_struct_exists { //If encrypted_struct doesn't exist in DataStore
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Do not have invitation struct/acess to file-> User is revoked")))
			return nil, errors.New(strings.ToTitle("Do not have invitation struct/acess to file-> User is revoked"))
		}
		//Decrypt invitation struct
		var invite, e10 = DecryptStruct(userdata, &encrypted_struct, invite_UUID, filename)
		if e10 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Decrypted invitation struct incorrectly.")))
			return nil, errors.New(strings.ToTitle("Decrypted invitation structincorrectly."))
		}

		/* 3. Get the encrypted file contents. Check that the HMAC tag saved in file _contents matches the (true)
		recomputed HMAC tag
		-> checks if file content has been tampered with while the user was gone
		-> IF tags match, decrypt file contents and return that
		*/
		//HELPER FUNCTION THAT DOES ALL THAT: DecryptContents(userdata *User, invite *Invitation, filename string) (content []byte, err error)
		file_contents, err3 := DecryptContents(userdata, invite)
		if err3 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Decrypted files incorrectly.")))
			return nil, errors.New(strings.ToTitle("Decrypted files incorrectly."))
		}

		/**4. If you got the appended_file_contents, that means there were APPENDED file contents, and Decrypt contents decrypted EACH
		FILE CONTENT STRUCT, put it all together, and sent you ONE appended_file_content []byte.
		-> Need to create a new encrypted_contents_struct and put appended_file_content in there
		-> Need to get rid of all file contents_structs (they're now all saved in appended_file_content []byte)
		-> Need to make number_of_appends 0 again
		*/
		//i.Get file struct
		//Get File struct
		var file_struct File
		file_struct_in_bytes, exist := userlib.DatastoreGet(invite.Uuid_to_file)
		e2 := userlib.Unmarshal(file_struct_in_bytes, &file_struct)

		if !exist {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Couldn't get file from DataStore.")))
			return nil, errors.New(strings.ToTitle("Couldn't get file from DataStore."))
		}
		if e2 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Unmarshalled file incorrectly.")))
			return nil, errors.New(strings.ToTitle("Unmarshalled file incorrectly."))
		}

		//ii. Check if you got appended_file_contents
		if file_struct.Number_of_appends > 0 {
			appended_file_content := file_contents

			//Encrypt the appended_file contents and then create a new encrypted_contents_struct
			//and put appended_file_content in there
			encrypted_appended_file_contents, encrypted_and_HMAC_file_contents, err4 := EncryptContents(userdata, invite, appended_file_content)
			if err4 != nil {
				userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Encrypted appended_file_content incorrectly.")))
				return nil, errors.New(strings.ToTitle("Encrypted appended_file_content incorrectly."))
			}

			appended_file_content_struct_uuid := userlib.UUIDNew()
			appended_file_content_struct := &Encrypted_Struct_or_Contents{
				encrypted_appended_file_contents,
				encrypted_and_HMAC_file_contents,
			}

			//Get rid of all file contents_structs (they're now all saved in appended_file_content []byte)
			for i := 0; i < file_struct.Number_of_appends; i++ {
				//Deterministically getting the UUID of the appended file_content_structs
				temp := []byte(file_struct.Original_file_content_struct_uuid.String() + string(rune(i)))
				delete_UUID, err := userlib.UUIDFromBytes(userlib.Hash(temp))
				if err != nil {
					userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not deterministically create appended_file_contents' UUIDs!")))
					return nil, errors.New(strings.ToTitle("Did not deterministically create appended_file_contents' UUIDs!"))
				}
				userlib.DatastoreDelete(delete_UUID)
			}

			//Save appended content as the entire file content now
			//(File.original_file_content_uuid = appended_content_uuid)
			//Make File.number_of_appends 0 again
			file_struct.Original_file_content_struct_uuid = appended_file_content_struct_uuid
			file_struct.Number_of_appends = 0

			//iii. Save appended_file_content_struct in DataStore, and save/update File in Datastore
			//Save appended_file_content_struct in DataStore
			appended_file_content_struct_in_bytes, err5 := userlib.Marshal(appended_file_content_struct)
			userlib.DatastoreSet(appended_file_content_struct_uuid, appended_file_content_struct_in_bytes)
			if err5 != nil {
				userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Marshalled appended_file_content_struct incorrectly.")))
				return nil, errors.New(strings.ToTitle("Marshalled appended_file_content_struct incorrectly."))
			}

			//save/update File in Datastore
			file_struct_in_bytes, err6 := userlib.Marshal(file_struct)
			userlib.DatastoreSet(invite.Uuid_to_file, file_struct_in_bytes)
			if err6 != nil {
				userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Marshalled file_struct incorrectly.")))
				return nil, errors.New(strings.ToTitle("Marshalled file_struct incorrectly."))
			}

			return appended_file_content, nil

		} else {
			//There were no appended file contents. You're good to go
			return file_contents, nil
		}

	}
}

func (userdata *User) CreateInvitation(filename string, recipientUsername string) (invitationPtr userlib.UUID, err error) {
	if userdata == nil {
		return userlib.UUID{}, errors.New(strings.ToTitle("User has not logged in or was not created"))
	}

	//IN THE CASE where user is operating on the same file WITH 2 DIFFERENT USERS/DEVICES, SYNC/UPDATE USER STRUCT/POINTER
	username_in_bytes := []byte(userdata.Username)
	user_uuid, _ := userlib.UUIDFromBytes(userlib.Hash(username_in_bytes))

	syncing_user, _ := userlib.DatastoreGet(user_uuid)
	userlib.Unmarshal(syncing_user, userdata)

	//1.  Check that the fileName exists in caller's namespace
	//Checking if hashed fileName exists in Dict_of_Files
	filename_in_bytes := []byte(filename)
	hashed_file := userlib.Hash(filename_in_bytes)
	var file_flow, exists = userdata.Dict_of_Files[hex.EncodeToString(hashed_file)]

	//Follow path to file. If path to file is missing, throw an error
	if !exists { //If filepath doesn't exist, then throw an error
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("File with that file name does not exist")))
		return userlib.UUID{}, errors.New(strings.ToTitle("File with that file name does not exist"))
	} else {
		//2. Get Bob's public key and encrypt the message. Get user's private key and DS the (encrypted) message.
		purpose1 := userlib.Hash([]byte("public key encryption"))
		purpose2 := userlib.Hash([]byte("public key verification"))

		//Getting Bob's public encryption key
		//Instead of using string([]byte), it errors/non-deterministic. Use this instead:
		// Takes a sequence of bytes and renders as hex
		//h := hex.EncodeToString([]byte("fubar"))
		temp2 := recipientUsername + " " + hex.EncodeToString(purpose1)
		recipients_PKEnc_key, exists1 := userlib.KeystoreGet(temp2)
		if !exists1 {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Public key does not exist.")))
			return userlib.UUID{}, errors.New(strings.ToTitle("Public key does not exist."))
		}

		//Getting user's/Alice's private signature key
		var private_signature_key userlib.DSSignKey
		encrypted_private_decryption_key_struct_uuid := userdata.Dict_of_keys[hex.EncodeToString(purpose1)]
		encrypted_private_signature_key_struct_uuid, exists2 := userdata.Dict_of_keys[hex.EncodeToString(purpose2)]

		if !exists2 {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Problem in CreateAndSaveKeys() helper function.")))
			return userlib.UUID{}, errors.New(strings.ToTitle("Problem in CreateAndSaveKeys() helper function."))
		}

		_, private_signature_key, e2 := DecrypteAndGetKeys(userdata, encrypted_private_decryption_key_struct_uuid, encrypted_private_signature_key_struct_uuid)
		if e2 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Problem in DecrypteAndGetKeys() helper function.")))
			return userlib.UUID{}, errors.New(strings.ToTitle("Problem in DecrypteAndGetKeys() helper function."))
		}

		/*3. If username = file's ownername (i.e. invite.Ownername), then :
		i. Create an NEW invitation struct, and send IT'S UUID (authentically encrypted!!!) to the user invited/recipient
		ii. Add the new invitation struct's UUID to the UserAccessDict
		iii. Add the name of the user invited/recipient to UsersSentHere (in invitation struct) <-ALREADY DONE IN HELPER FUNCTION
		*/
		var inviteUUID userlib.UUID
		if file_flow.OwnerName == userdata.Username {
			//i. Create an NEW invitation struct, and send IT'S UUID (authentically encrypted!!!) to the user invited/recipient
			//_: new_encrypted_inv_struct
			_, uuid_invite, e3 := CreateInvitationStructToAnotherUser(userdata, filename, &file_flow, recipientUsername)
			inviteUUID = uuid_invite
			if e3 != nil {
				userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not create new (encrypted) invitaiton struct to a new user correctly.")))
				return userlib.UUID{}, errors.New(strings.ToTitle("Did not create new (encrypted) invitaiton struct to a new user correctly."))
			}

		} else {

			/*4. If username !- file's ownername, then:
			i. Use THE SAME INVITATION STRUCT and share that uuid to the user invited/recipient
			ii. Add user invited/recipient to UsersSentHere (in invitation struct)
			iii. Save changes made in the invitation's UsersSentHere
			*/
			//i. Share the SAME INVITATION STRUCT with user
			inviteUUID = file_flow.Uuid_to_invitation_struct
			//ii. Add user invited/recipient to UsersSentHere (in invitation struct)
			var encrypted_inv_struct Encrypted_Struct_or_Contents
			encrypted_invite_struct_in_bytes, exists4 := userlib.DatastoreGet(inviteUUID)
			userlib.Unmarshal(encrypted_invite_struct_in_bytes, &encrypted_inv_struct)
			if !exists4 {
				userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Unmarshalled incorrectly.")))
				return userlib.UUID{}, errors.New(strings.ToTitle("Unmarshalled incorrectly."))
			}

			invite_struct, e4 := DecryptStruct(userdata, &encrypted_inv_struct, inviteUUID, filename)
			if e4 != nil {
				userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Decrypted invitation struct incorrectly.")))
				return userlib.UUID{}, errors.New(strings.ToTitle("Decrypted invitation struct incorrectly."))
			}

			list_of_people := invite_struct.Users_sent_here
			added_list_of_people := append(list_of_people, recipientUsername)
			invite_struct.Users_sent_here = added_list_of_people

			//iii. Save changes made in invite struct
			encrypted_struct, encrpyt_and_HMAC_inv_struct, err := EncryptInvitationStruct(userdata, invite_struct, inviteUUID)
			if err != nil {
				userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not encrypt-then_HMAC invitation struct correctly!")))
				return userlib.UUID{}, errors.New(strings.ToTitle("Did not encrypt-then_HMAC correctly!"))
			}

			auth_enc_struct := &Encrypted_Struct_or_Contents{
				encrypted_struct,
				encrpyt_and_HMAC_inv_struct,
			}

			//Saving to DataStore
			auth_enc_inv_struct_bytes, e := userlib.Marshal(auth_enc_struct)
			userlib.DatastoreSet(inviteUUID, auth_enc_inv_struct_bytes)
			if e != nil {
				return userlib.UUID{}, errors.New(strings.ToTitle("Paniiiccccc"))
			}
		}

		//5. Once you have the new/old invitation struct UUID you want to send/share with user, you must
		//SEND AN INVITATION MESSAGE WITH INVITATION UUID INSIDE,
		//PUBLICLY ENCRYPT AND DS THE INVITATION MESSAGE so user who got it will be able to verify & use it

		//Need BOB'S PUBLIC KEY to encrypt and ALICE'S PRIVATE KEY to sign

		//InviteUUID must be in bytes -> Marshal!
		inviteUUID_in_bytes, e5 := userlib.Marshal(inviteUUID)
		sending_invite_uuid_using_PKE, e6 := userlib.PKEEnc(recipients_PKEnc_key, inviteUUID_in_bytes)
		sending_invite_uuid_using_DS_tag, e7 := userlib.DSSign(private_signature_key, sending_invite_uuid_using_PKE)

		if e5 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Marshalled incorrectly.")))
			return userlib.UUID{}, errors.New(strings.ToTitle("Marshalled incorrectly."))
		}
		if e6 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Publicly encrypted incorrectly.")))
			return userlib.UUID{}, errors.New(strings.ToTitle("Publicly encrypted incorrectly."))
		}
		if e7 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Publicly signed (DS) incorrectly.")))
			return userlib.UUID{}, errors.New(strings.ToTitle("Publicly signed (DS) incorrectly."))
		}

		//Put publibly encrypted invite message all in a struct
		invite_message_struct := &Publicly_Encrypted_Contents{
			sending_invite_uuid_using_PKE,    //Ciphertext(invite_message)
			sending_invite_uuid_using_DS_tag, //DS_tag of (ciphertext)
		}

		//Save that in the Datastore, and send its key/UUID to the reciepient user
		invite_message_UUID := userlib.UUIDNew()
		invite_message_struct_in_bytes, e8 := userlib.Marshal(invite_message_struct)
		userlib.DatastoreSet(invite_message_UUID, invite_message_struct_in_bytes)
		if e8 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Saved in Datastore incorrectly.")))
			return userlib.UUID{}, errors.New(strings.ToTitle("Saved in Datastore incorrectly."))
		}

		//Saave useraccessdict
		return invite_message_UUID, nil
	}

}

//Helper function
//Returns the private decryption and signature keys
func DecrypteAndGetKeys(userdata *User, encrypted_private_decryption_key_struct_uuid userlib.UUID, encrypted_private_signature_key_struct_uuid userlib.UUID) (private_decryption_key userlib.PKEDecKey, private_signature_key userlib.DSSignKey, err error) {

	//1. Get the symmetric_encryption_key and HMAC_key
	PKey_UUID_in_bytes, e5 := userlib.Marshal(encrypted_private_decryption_key_struct_uuid)
	incremented_PK_UUID := userlib.Hash([]byte(hex.EncodeToString(PKey_UUID_in_bytes) + "100"))
	PK_symmetric_encryption_key := userlib.Argon2Key(PKey_UUID_in_bytes, incremented_PK_UUID, 16)
	PK_HMAC_key, e6 := userlib.HashKDF(PKey_UUID_in_bytes[:16], []byte("PK_private_key's HMAC key"))

	DSKey_UUID_in_bytes, e7 := userlib.Marshal(encrypted_private_signature_key_struct_uuid)
	incremented_DSKey_UUID := userlib.Hash([]byte(hex.EncodeToString(DSKey_UUID_in_bytes) + "100"))
	DS_symmetric_encryption_key := userlib.Argon2Key(DSKey_UUID_in_bytes, incremented_DSKey_UUID, 16)
	DS_HMAC_key, e8 := userlib.HashKDF(DSKey_UUID_in_bytes[:16], []byte("DS_private_key's HMAC key"))

	if e5 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("PK_private_key_struct's UUID marshalled incorrectly.")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("PK_private_key_struct's UUID marshalled incorrectly."))
	}
	if e7 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("DS_private_key_struct's UUID marshalled incorrectly.")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("DS_private_key_struct's UUID marshalled incorrectly."))
	}
	if e6 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("PK_private_key's HMAC made incorrectly.")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("PK_private_key's HMAC made incorrectly."))
	}
	if e8 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("DS_private_key's HMAC key made correctly.")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("DS_private_key's HMAC key made correctly."))
	}

	//2. Get the encrypted private keys' struct
	var encrypted_private_decryption_key_struct Encrypted_Struct_or_Contents
	encrypted_private_decryption_key_struct_in_bytes, exists1 := userlib.DatastoreGet(encrypted_private_decryption_key_struct_uuid)
	userlib.Unmarshal(encrypted_private_decryption_key_struct_in_bytes, &encrypted_private_decryption_key_struct)

	var encrypted_private_signature_key_struct Encrypted_Struct_or_Contents
	encrypted_private_signature_key_struct_in_bytes, exists2 := userlib.DatastoreGet(encrypted_private_signature_key_struct_uuid)
	userlib.Unmarshal(encrypted_private_signature_key_struct_in_bytes, &encrypted_private_signature_key_struct)

	if !exists1 {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Cannot get the encrypted_PK_private_key from Datastore.")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("Cannot get the encrypted_PK_private_key from Datastore."))
	}
	if !exists2 {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Cannot get the encrypted_DS_private_key from Datastore.")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("Cannot get the encrypted_DS_private_key from Datastore."))
	}

	//3. Make sure the HMAC_tags match

	//Check HMAC_tag matches encrypted struct's HMAC_tag
	//i. For PKey:
	cipher_text := encrypted_private_decryption_key_struct.Ciphertext
	tester_tag := encrypted_private_decryption_key_struct.HMAC_tag
	real_tag, err := userlib.HMACEval(PK_HMAC_key[:16], cipher_text)
	if err != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not encrypt-then_HMAC correctly!")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("Did not encrypt-then_HMAC correctly!"))
	}
	//Checking if they're equal
	//HMACEqual(a []byte, b []byte) (bool)
	the_same := userlib.HMACEqual(tester_tag, real_tag)
	if !the_same {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("HMAC_tags don't match!!!")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("HMAC_tags don't match!!!"))
	}

	//ii. For DSKey:
	cipher_text2 := encrypted_private_signature_key_struct.Ciphertext
	tester_tag2 := encrypted_private_signature_key_struct.HMAC_tag
	real_tag2, err := userlib.HMACEval(DS_HMAC_key[:16], cipher_text2)
	if err != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not encrypt-then_HMAC correctly!")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("Did not encrypt-then_HMAC correctly!"))
	}
	//Checking if they're equal
	//HMACEqual(a []byte, b []byte) (bool)
	the_same2 := userlib.HMACEqual(tester_tag2, real_tag2)
	if !the_same2 {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("HMAC_tags don't match!!!")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("HMAC_tags don't match!!!"))
	}

	//4. If they do match, decrypt the ciphertext in encrypted_private_keys struct to get the private keys
	var PK_private_key userlib.PKEDecKey
	private_decryption_key_in_bytes := userlib.SymDec(PK_symmetric_encryption_key, cipher_text)
	err2 := userlib.Unmarshal(private_decryption_key_in_bytes, &PK_private_key)

	var DS_private_key userlib.DSSignKey
	private_verification_key_in_bytes := userlib.SymDec(DS_symmetric_encryption_key, cipher_text2)
	err3 := userlib.Unmarshal(private_verification_key_in_bytes, &DS_private_key)

	if err2 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Unmarshalled incorrectly.")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("Unmarshalled incorrectly."))
	}
	if err3 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Unmarshalled incorrectly.")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("Unmarshalled incorrectly."))
	}

	//5. Return the private keys

	return PK_private_key, DS_private_key, nil

}

//HelperFunction: Create and save public encryption keys and DS keys (both private and public)
//Returns UUID of the encrypted private keys_structs
func CreateAndSaveKeys(userdata *User) (private_decryption_key userlib.PKEDecKey, private_signature_key userlib.DSSignKey, err error) {
	purpose1 := userlib.Hash([]byte("public key encryption"))
	purpose2 := userlib.Hash([]byte("public key verification"))

	//1. Generate the keys
	//PKEKeyGen() (PKEEncKey, PKEDecKey, error)
	//DSKeyGen() (DSSignKey, DSVerifyKey, error)
	public_encryption_key, private_decryption_key, err1 := userlib.PKEKeyGen()
	private_signature_key, public_verification_key, err2 := userlib.DSKeyGen()
	if err1 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("PKEKeyGen wasn't used correctly.")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("PKEKeyGen wasn't used correctly."))
	}
	if err2 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("DSKeyGen wasn't used correctly.")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("DSKeyGen wasn't used correctlyz"))
	}

	//2. Save in Keystore: KeystoreSet(key string, value PKEEncKey/DSVerifyKey) (error)
	//Instead of using string([]byte), it errors/non-deterministic. Use this instead:
	//Takes a sequence of bytes and renders as hex
	//h := hex.EncodeToString([]byte("fubar"))
	temp1 := userdata.Username + " " + hex.EncodeToString(purpose1)
	err3 := userlib.KeystoreSet(temp1, public_encryption_key)

	temp2 := userdata.Username + " " + hex.EncodeToString(purpose2)
	err4 := userlib.KeystoreSet(temp2, public_verification_key)

	if err3 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Public encryption key was not stored in KeyStore correctly.")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("Public encryption key was not stored in KeyStore correctly."))
	}
	if err4 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Public verification key was not stored in KeyStore correctly.")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("Public verification key was not stored in KeyStore correctly."))
	}

	//3. Save in Datastore:
	//i. Create new_encrypted_structs_UUID for keys
	PKey_UUID := userlib.UUIDNew()
	DSKey_UUID := userlib.UUIDNew()

	//ii. Generate symmetric keys & HMACs to authenically encrypt each private key using the newUUID
	//Creating HMAC and encrypted keys using encrypted_private_keys struct's UUID
	PKey_UUID_in_bytes, e5 := userlib.Marshal(PKey_UUID)
	incremented_PK_UUID := userlib.Hash([]byte(hex.EncodeToString(PKey_UUID_in_bytes) + "100"))
	PK_symmetric_encryption_key := userlib.Argon2Key(PKey_UUID_in_bytes, incremented_PK_UUID, 16)
	PK_HMAC_key, e6 := userlib.HashKDF(PKey_UUID_in_bytes[:16], []byte("PK_private_key's HMAC key"))

	DSKey_UUID_in_bytes, e7 := userlib.Marshal(DSKey_UUID)
	incremented_DSKey_UUID := userlib.Hash([]byte(hex.EncodeToString(DSKey_UUID_in_bytes) + "100"))
	DS_symmetric_encryption_key := userlib.Argon2Key(DSKey_UUID_in_bytes, incremented_DSKey_UUID, 16)
	DS_HMAC_key, e8 := userlib.HashKDF(DSKey_UUID_in_bytes[:16], []byte("DS_private_key's HMAC key"))

	if e5 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("PK_private_key_struct's UUID marshalled incorrectly.")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("PK_private_key_struct's UUID marshalled incorrectly."))
	}
	if e7 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("DS_private_key_struct's UUID marshalled incorrectly.")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("DS_private_key_struct's UUID marshalled incorrectly."))
	}
	if e6 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("PK_private_key's HMAC made incorrectly.")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("PK_private_key's HMAC made incorrectly."))
	}
	if e8 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("DS_private_key's HMAC key made correctly.")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("DS_private_key's HMAC key made correctly."))
	}

	//iii. Auth. Encrypt private keys
	iv1 := userlib.RandomBytes(16)
	iv2 := userlib.RandomBytes(16)
	private_decryption_key_in_bytes, e51 := userlib.Marshal(private_decryption_key)
	private_signature_key_in_bytes, e50 := userlib.Marshal(private_signature_key)
	if e51 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Marshalled incorrectly.")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("Marshalled incorrectly."))
	}
	if e50 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Marshalled incorrectly.")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("Marshalled incorrectly."))
	}

	PK_ciphertext := userlib.SymEnc(PK_symmetric_encryption_key, iv1, private_decryption_key_in_bytes)
	PK_HMAC_tag, e9 := userlib.HMACEval(PK_HMAC_key[:16], PK_ciphertext)

	DS_ciphertext := userlib.SymEnc(DS_symmetric_encryption_key, iv2, private_signature_key_in_bytes)
	DS_HMAC_tag, e10 := userlib.HMACEval(DS_HMAC_key[:16], DS_ciphertext)

	if e9 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("PK_HMAC_tag made correctly.")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("PK_HMAC_tag made correctly."))
	}
	if e10 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("DS_HMAC_tag made correctly.")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("DS_HMAC_tag made correctly."))
	}

	//iv. Create new_encrypted_structs saving the authentically encrypted private key
	encrypted_PK_private_key_struct := &Encrypted_Struct_or_Contents{
		PK_ciphertext,
		PK_HMAC_tag,
	}

	encrypted_DS_private_key_struct := &Encrypted_Struct_or_Contents{
		DS_ciphertext,
		DS_HMAC_tag,
	}

	//v. Save those new_encrypted_structs in DataStore, where key = newUUID, and Marshal(new_encrypted_structs)

	encrypted_PK_private_key_struct_in_bytes, e11 := userlib.Marshal(encrypted_PK_private_key_struct)
	userlib.DatastoreSet(PKey_UUID, encrypted_PK_private_key_struct_in_bytes)
	if e11 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not save encrypted_PK_private_key_struct correctly.")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("Did not save encrypted_PK_private_key_struct correctly."))
	}

	encrypted_DS_private_key_struct_in_bytes, e12 := userlib.Marshal(encrypted_DS_private_key_struct)
	userlib.DatastoreSet(DSKey_UUID, encrypted_DS_private_key_struct_in_bytes)
	if e12 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not save encrypted_DS_private_key_struct correctly.")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("Did not save encrypted_DS_private_key_struct correctly."))
	}

	//vi. Save those new_encrypted_struct_uuids in Dict_of_keys, where key = hashed(purpose) and value = new_encrypted_structs_uuids
	//INSTEAD OF SAVING IT IN DATASTORE, WHY NOT SAVE IT DIRECTLY IN DICT_OF_KEYS MAP????
	userdata.Dict_of_keys[hex.EncodeToString(purpose1)] = PKey_UUID
	userdata.Dict_of_keys[hex.EncodeToString(purpose2)] = DSKey_UUID

	//vii. Save the updated Dict_of_keys, by saving the User (again) to the Datastore
	username_in_bytes := []byte(userdata.Username)
	Uuid_of_user, e1 := userlib.UUIDFromBytes(userlib.Hash(username_in_bytes))
	user_struct_in_bytes, e2 := userlib.Marshal(userdata)
	userlib.DatastoreSet(Uuid_of_user, user_struct_in_bytes)

	if e1 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not get user's UUID correctly.")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("Did not get user's UUID correctly."))
	}
	if e2 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not marshal user struct correctly.")))
		return userlib.PKEDecKey{}, userlib.DSSignKey{}, errors.New(strings.ToTitle("Did not marshal user struct correctly."))
	}

	//Return the new private keys
	return private_decryption_key, private_signature_key, nil
}

//Helper Function: If it's the owner, they must create a NEW INVITATION STRUCT and send it to the user recipient
//input: FilePath/Vertex pointer
//output: encrypted invitation struct and its UUID
func CreateInvitationStructToAnotherUser(userdata *User, filename string, file_flow *Vertex, new_user_adding string) (encrypted_invitation_struct *Encrypted_Struct_or_Contents, inviteUUID userlib.UUID, err error) {
	//Retrieve the file you want to share
	owners_invite_to_file_uuid := file_flow.Uuid_to_invitation_struct
	owners_invite_to_file_in_bytes, exists3 := userlib.DatastoreGet(owners_invite_to_file_uuid)

	var encrypted_owners_invitation_struct_to_file Encrypted_Struct_or_Contents
	err2 := userlib.Unmarshal(owners_invite_to_file_in_bytes, &encrypted_owners_invitation_struct_to_file)
	owners_invite_to_file, err3 := DecryptStruct(userdata, &encrypted_owners_invitation_struct_to_file, owners_invite_to_file_uuid, filename)

	fileUUID := owners_invite_to_file.Uuid_to_file

	if !exists3 {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not get owner's (encrypted) invitation struct from Datastore correctly.")))
		return &Encrypted_Struct_or_Contents{}, userlib.UUID{}, errors.New(strings.ToTitle("Did not get owner's (encrypted) invitation struct from Datastore correctly."))
	}
	if err2 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Marshalled incorrectly.")))
		return &Encrypted_Struct_or_Contents{}, userlib.UUID{}, errors.New(strings.ToTitle("Marshalled incorrectly."))
	}
	if err3 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Decrypted owner's invitation struct incorrectly.")))
		return &Encrypted_Struct_or_Contents{}, userlib.UUID{}, errors.New(strings.ToTitle("Decrypted owner's invitation struct incorrectly."))
	}

	//Creating invitation struct to send to user
	invitation_struct_UUID := userlib.UUIDNew()

	var users_sent_here []string
	users_sent_here = append(users_sent_here, new_user_adding)

	invite := &Invitation{ //Will be encrypted with symmetric key = HKDF{UUID to invitation struct}
		file_flow.OwnerName,
		users_sent_here,
		fileUUID, //Sharing owner's file with user
	}

	//Encrypt entire Invitation struct so only non-malicious user can see it
	encrypted_struct, encrpyt_and_HMAC_inv_struct, err := EncryptInvitationStruct(userdata, invite, invitation_struct_UUID)
	if err != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not encrypt-then_HMAC invitation struct correctly!")))
		return &Encrypted_Struct_or_Contents{}, userlib.UUID{}, errors.New(strings.ToTitle("Did not encrypt-then_HMAC correctly!"))
	}

	//Auh_enc_inv_struct = { ENC_k1(inv_struct_in_bytes), MAC_k2(ENC_k1(inv_struct_in_bytes))}
	auth_enc_struct := &Encrypted_Struct_or_Contents{
		encrypted_struct,
		encrpyt_and_HMAC_inv_struct,
	}

	//ii. Save to DataStore
	auth_enc_inv_struct_bytes, e := userlib.Marshal(auth_enc_struct)
	userlib.DatastoreSet(invitation_struct_UUID, auth_enc_inv_struct_bytes)
	if e != nil {
		return &Encrypted_Struct_or_Contents{}, userlib.UUID{}, errors.New(strings.ToTitle("Paniiiccccc"))
	}

	//ii. Add the new invitation struct's UUID to owner's UserAccessDict
	list_of_names := file_flow.UserAccessDict[userdata.Username]
	added_list_of_names := append(list_of_names, invitation_struct_UUID)
	file_flow.UserAccessDict[userdata.Username] = added_list_of_names

	//iii. Save/Update UserAccessDict in Datastore
	username_in_bytes := []byte(userdata.Username)
	Uuid_of_user, e1 := userlib.UUIDFromBytes(userlib.Hash(username_in_bytes))
	user_struct_in_bytes, e2 := userlib.Marshal(userdata)
	userlib.DatastoreSet(Uuid_of_user, user_struct_in_bytes)

	if e1 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not get user's UUID correctly.")))
		return &Encrypted_Struct_or_Contents{}, userlib.UUID{}, errors.New(strings.ToTitle("Did not get user's UUID correctly."))
	}
	if e2 != nil {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not marshal user struct correctly.")))
		return &Encrypted_Struct_or_Contents{}, userlib.UUID{}, errors.New(strings.ToTitle("Did not marshal user struct correctly."))
	}

	return auth_enc_struct, invitation_struct_UUID, nil
}

func (userdata *User) AcceptInvitation(senderUsername string, invitationPtr userlib.UUID, filename string) error {
	if userdata == nil {
		return errors.New(strings.ToTitle("User has not logged in or was not created"))
	}

	//IN THE CASE where user is operating on the same file WITH 2 DIFFERENT USERS/DEVICES, SYNC/UPDATE USER STRUCT/POINTER
	username_in_bytes := []byte(userdata.Username)
	user_uuid, _ := userlib.UUIDFromBytes(userlib.Hash(username_in_bytes))

	syncing_user, _ := userlib.DatastoreGet(user_uuid)
	userlib.Unmarshal(syncing_user, userdata)

	//1. Check if filename is in user's namespace. If it IS, throw an error!
	//Checking if hashed fileName exists in Dict_of_Files
	filename_in_bytes := []byte(filename)
	hashed_file := userlib.Hash(filename_in_bytes)
	_, exists := userdata.Dict_of_Files[hex.EncodeToString(hashed_file)]

	//If filepath does exist, then throw an error
	if exists {
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("File with that file name does not exist")))
		return errors.New(strings.ToTitle("File with that file name does not exist"))
	} else {
		//2. Grab user's (getting invitation) PRIVATE DECRYPTION KEY and sender's PUBLIC KEY (for verification)
		purpose1 := userlib.Hash([]byte("public key encryption"))
		purpose2 := userlib.Hash([]byte("public key verification"))

		//Getting sender's public key verification (verify Alice sent message to Bob)
		temp := senderUsername + " " + hex.EncodeToString(purpose2)
		public_verification_key, exists := userlib.KeystoreGet(temp)
		if !exists {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Public keys were not created for this user in InitUser()")))
			return errors.New(strings.ToTitle("Public keys were not created for this user in InitUser()"))
		}

		//Getting user's private decryption key (Bob can decrypt message Alice sent, cuz Alic used Bob's PUBLIC enc key)
		//FLAW! OWNER NAME IS SAVED NO WHERE HERE! VERTEX STRUCT OF USER RECIPIENT HASNT BEEN MADE IT, AND INVITE STRUCT IS ENCRYPTED!
		var private_decryption_key userlib.DSSignKey
		encrypted_private_decryption_key_struct_uuid := userdata.Dict_of_keys[hex.EncodeToString(purpose1)]
		encrypted_private_signature_key_struct_uuid, exists2 := userdata.Dict_of_keys[hex.EncodeToString(purpose2)]

		if !exists2 {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Problem in CreateAndSaveKeys() helper function.")))
			return errors.New(strings.ToTitle("Problem in CreateAndSaveKeys() helper function."))
		}

		private_decryption_key, _, e2 := DecrypteAndGetKeys(userdata, encrypted_private_decryption_key_struct_uuid, encrypted_private_signature_key_struct_uuid)
		if e2 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Problem in DecrypteAndGetKeys() helper function.")))
			return errors.New(strings.ToTitle("Problem in DecrypteAndGetKeys() helper function."))
		}

		//3. Check out invitationPtr uuid -> it's an Encrypted_Struct type, where the decrypted(ciphertext) is the invitation struct's uuid
		//Check out the invitationPtr uuid/ciphertext/M
		var encrypted_invite_message_struct Publicly_Encrypted_Contents
		encrypted_invite_message_struct_in_bytes, exists3 := userlib.DatastoreGet(invitationPtr)
		e3 := userlib.Unmarshal(encrypted_invite_message_struct_in_bytes, &encrypted_invite_message_struct)
		if !exists3 {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("invite_message not saved in DataStore correctly.")))
			return errors.New(strings.ToTitle("invite_message not saved in DataStore correctly."))
		}
		if e3 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Unmarshalled incorrectly.")))
			return errors.New(strings.ToTitle("Unmarshalled incorrectly."))
		}

		//Verifying
		cipher_text := encrypted_invite_message_struct.Ciphertext
		DS_tag := encrypted_invite_message_struct.DS_tag

		e4 := userlib.DSVerify(public_verification_key, cipher_text, DS_tag)
		if e4 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not DSVerify correctly.")))
			return errors.New(strings.ToTitle("Did not DSVerify correctly."))
		}

		//Decrypting (to get the ACTUAL invite uuid, which leads to invite struct)
		var invite_struct_uuid userlib.UUID
		invite_struct_uuid_in_bytes, e5 := userlib.PKEDec(private_decryption_key, cipher_text)
		e6 := userlib.Unmarshal(invite_struct_uuid_in_bytes, &invite_struct_uuid)
		if e5 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not PKDecrypt correctly.")))
			return errors.New(strings.ToTitle("Did not PKDecrypt correctly."))
		}
		if e6 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Unmarshalled incorrectly.")))
			return errors.New(strings.ToTitle("Unmarshalled incorrectly."))
		}

		//4. Decrypt invitation struct, and create Vertex struct to link inv. struct to user
		//Getting encrypted struct
		var encrypted_inv_struct Encrypted_Struct_or_Contents
		encrypted_invite_struct_in_bytes, exists5 := userlib.DatastoreGet(invite_struct_uuid)
		e7 := userlib.Unmarshal(encrypted_invite_struct_in_bytes, &encrypted_inv_struct)
		if !exists5 {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not get encrypted_inv_struct from DataStore correctly.")))
			return errors.New(strings.ToTitle("Did not get encrypted_inv_struct from DataStore correctly."))
		}
		if e7 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Unmarshalled incorrectly.")))
			return errors.New(strings.ToTitle("Unmarshalled incorrectly."))
		}

		//Getting (decrypted) invitation struct
		var invite_struct *Invitation
		invite_struct, e8 := DecryptStruct(userdata, &encrypted_inv_struct, invite_struct_uuid, filename)
		if e8 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not decrypt invitation struct correctly.")))
			return errors.New(strings.ToTitle("Did not decrypt invitation struct correctly."))
		}

		//Create Vertex struct to link inv. struct to User
		file_flow := &Vertex{
			invite_struct.OwnerName,
			invite_struct_uuid,
			nil,
		}

		//5. Link Vertex to User struct by adding it to Dict_of_Files, where the key is the hashed filename
		userdata.Dict_of_Files[hex.EncodeToString(hashed_file)] = *file_flow
		//6. SAVE ALL CHANGES
		//i. Save UserStruct changes (which saves the addition of Vertex in Dic_of_files)
		username_in_bytes := []byte(userdata.Username)
		Uuid_of_user, e9 := userlib.UUIDFromBytes(userlib.Hash(username_in_bytes))
		userdata_in_bytes, e10 := userlib.Marshal(userdata)
		userlib.DatastoreSet(Uuid_of_user, userdata_in_bytes)

		if e9 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Created Uuid_of_user incorrectly.")))
			return errors.New(strings.ToTitle("Created Uuid_of_user incorrectly."))
		}
		if e10 != nil {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Marshalled incorrectly.")))
			return errors.New(strings.ToTitle("Marshalled incorrectly."))
		}

		return nil
	}
}

//TODO: Store shared file (using mostly of the same code from storefile)

func (userdata *User) RevokeAccess(filename string, recipientUsername string) error {
	if userdata == nil {
		return errors.New(strings.ToTitle("User has not logged in or was not created"))
	}

	//IN THE CASE where user is operating on the same file WITH 2 DIFFERENT USERS/DEVICES, SYNC/UPDATE USER STRUCT/POINTER
	username_in_bytes := []byte(userdata.Username)
	user_uuid, _ := userlib.UUIDFromBytes(userlib.Hash(username_in_bytes))

	syncing_user, _ := userlib.DatastoreGet(user_uuid)
	userlib.Unmarshal(syncing_user, userdata)

	//1.  Check that the fileName exists in caller's namespace
	//Checking if hashed fileName exists in Dict_of_Files
	filename_in_bytes := []byte(filename)
	hashed_file := userlib.Hash(filename_in_bytes)
	var file_flow, exists = userdata.Dict_of_Files[hex.EncodeToString(hashed_file)]

	//Follow path to file. If path to file is missing, throw an error
	if !exists { //If filepath doesn't exist, then throw an error
		userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("File with that file name does not exist")))
		return errors.New(strings.ToTitle("File with that file name does not exist"))
	} else {
		// 2. Check if the user (revoking access) is the OWNER
		if file_flow.OwnerName != userdata.Username {
			userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Only the owner of the file can revoke access.")))
			return errors.New(strings.ToTitle("Only the owner of the file can revoke access."))
		} else {
			// 3. Change the file's location/UUID, and save it in OWNER'S invitation struct
			//i. Get the owner's invitation struct and that will give us file's current location/UUID
			var encrypted_inv_struct Encrypted_Struct_or_Contents
			encrypted_inv_struct_in_bytes, exists2 := userlib.DatastoreGet(file_flow.Uuid_to_invitation_struct)
			e1 := userlib.Unmarshal(encrypted_inv_struct_in_bytes, &encrypted_inv_struct)
			if !exists2 {
				userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Invitation struct is not in the DataStore.")))
				return errors.New(strings.ToTitle("Invitation struct is not in the DataStore."))
			}
			if e1 != nil {
				userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Unmarshalled incorrectly.")))
				return errors.New(strings.ToTitle("Unmarshalled incorrectly."))
			}

			var inv_struct *Invitation
			inv_struct, e2 := DecryptStruct(userdata, &encrypted_inv_struct, file_flow.Uuid_to_invitation_struct, filename)
			if e2 != nil {
				userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Decrypted struct incorrectly.")))
				return errors.New(strings.ToTitle("Decrypted struct incorrectly."))
			}

			//ii. Update the owner's invitation struct so that it 's file's location is the new_file_uuid

			old_fileUUID := inv_struct.Uuid_to_file
			//Access old file struct and get its encrypted file contents
			file_contents, e22 := DecryptContents(userdata, inv_struct)
			if e22 != nil {
				userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Decrypted file contents incorrectly.")))
				return errors.New(strings.ToTitle("Decrypted file contents incorrectly."))
			}

			//Re-encrypt the file contents using new fileUUID as key
			new_file_uuid := userlib.UUIDNew()
			new_file_content_struct_UUID := userlib.UUIDNew()

			inv_struct.Uuid_to_file = new_file_uuid
			new_encrypted_content, new_encrpyt_and_HMAC_content, e23 := EncryptContents(userdata, inv_struct, file_contents)
			if e23 != nil {
				userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Enrypted file contents incorrectly.")))
				return errors.New(strings.ToTitle("Encrypted file contents incorrectly."))
			}

			//Save newly-encrypted contents in new file_strct struct
			//Saving content into a new file_content_struct.
			new_file_content_struct := &Encrypted_Struct_or_Contents{
				new_encrypted_content,
				new_encrpyt_and_HMAC_content,
			}

			new_file_struct := &File{
				new_file_content_struct_UUID,
				0,
			}

			//Save new file contents in Datastore, and new file in Datastore
			new_file_content_struct_in_bytes, e24 := userlib.Marshal(new_file_content_struct)
			userlib.DatastoreSet(new_file_content_struct_UUID, new_file_content_struct_in_bytes)
			if e24 != nil {
				userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Marshalled file content struct incorrectly.")))
				return errors.New(strings.ToTitle("Marshalled file contents incorrectly."))
			}

			new_file_struct_in_bytes, e25 := userlib.Marshal(new_file_struct)
			userlib.DatastoreSet(new_file_uuid, new_file_struct_in_bytes)
			if e25 != nil {
				userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Marshalled file struct incorrectly.")))
				return errors.New(strings.ToTitle("Marshalled file content struct incorrectly."))
			}

			//iii. Delete file's old location from Datastore
			userlib.DatastoreDelete(old_fileUUID)

			/* 4. Go through owner's file's AccessDict, and LOOP THRU each inv_struct_uuid
			i. Access the DECRYPTED invitation struct
			ii. In the inv_struct, UPDATE EACH STRUCT'S FILEUUID
				-> HOWEVER, IF the reciepientUsername IS IN the Users_sent_here dictionary,
					DON'T UPDATE THE FILE_UUID
			*/
			var new_list []userlib.UUID
			for i := 0; i < len(file_flow.UserAccessDict[userdata.Username]); i++ {

				shared_users_invite_uuid := file_flow.UserAccessDict[userdata.Username][i]

				//i. Access the DECRYPTED invitation struct
				var shared_users_encrypted_invite *Encrypted_Struct_or_Contents
				shared_users_encrypted_invite_in_bytes, exists := userlib.DatastoreGet(shared_users_invite_uuid)
				err := userlib.Unmarshal(shared_users_encrypted_invite_in_bytes, &shared_users_encrypted_invite)
				if !exists {
					userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Shared encrypted struct is not in the DataStore.")))
					return errors.New(strings.ToTitle("Shared encrypted struct is not in the DataStore."))
				}
				if err != nil {
					userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Unmarshalled incorrectly.")))
					return errors.New(strings.ToTitle("Unmarshalled incorrectly."))
				}

				shared_users_invite, e2 := DecryptStruct(userdata, shared_users_encrypted_invite, shared_users_invite_uuid, filename)
				if e2 != nil {
					userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Decrypted struct incorrectly.")))
					return errors.New(strings.ToTitle("Decrypted struct incorrectly."))
				}

				/*ii. In the inv_struct, UPDATE EACH STRUCT'S FILEUUID
				-> HOWEVER, IF the reciepientUsername IS IN the Users_sent_here string list,
					DON'T UPDATE THE FILE_UUID
				*/
				users_sent_here_list := shared_users_invite.Users_sent_here
				recipients_is_here := false
				for k := 0; k < len(users_sent_here_list); k++ {
					users_sent_here := users_sent_here_list[k]
					//IF the reciepientUsername IS IN the Users_sent_here string list,
					//DON'T UPDATE THE FILE_UUID
					if users_sent_here == recipientUsername {
						recipients_is_here = true
					}
				}
				if recipients_is_here {
					continue
				}

				new_list = append(new_list, shared_users_invite_uuid)
				shared_users_invite.Uuid_to_file = new_file_uuid

				//iii. Save each struct's changes in Datastore
				encrypted_shared_struct, encrpyt_and_HMAC_shared_inv_struct, err := EncryptInvitationStruct(userdata, shared_users_invite, shared_users_invite_uuid)
				if err != nil {
					userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not encrypt-then_HMAC invitation struct correctly!")))
					return errors.New(strings.ToTitle("Did not encrypt-then_HMAC correctly!"))
				}

				auth_enc_struct := &Encrypted_Struct_or_Contents{
					encrypted_shared_struct,
					encrpyt_and_HMAC_shared_inv_struct,
				}

				//Saving to DataStore
				auth_enc_shared_inv_struct_bytes, e := userlib.Marshal(auth_enc_struct)
				userlib.DatastoreSet(shared_users_invite_uuid, auth_enc_shared_inv_struct_bytes)
				if e != nil {
					return errors.New(strings.ToTitle("Paniiiccccc"))
				}
			}

			//Update UserAccessDict so that it doesn't have recipient in it
			file_flow.UserAccessDict[userdata.Username] = new_list

			//5. SAVE ALL CHANGES
			//i. Owner's UserAccessDict changed -> save/update owner's user struct in Datastore
			username_in_bytes := []byte(userdata.Username)
			uuid_of_user, e10 := userlib.UUIDFromBytes(userlib.Hash(username_in_bytes))
			user_in_bytes, e11 := userlib.Marshal(userdata)
			userlib.DatastoreSet(uuid_of_user, user_in_bytes)
			if e10 != nil {
				userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("panic??")))
				return errors.New(strings.ToTitle("panic??"))
			}
			if e11 != nil {
				userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Marshalled incorrectly!")))
				return errors.New(strings.ToTitle("Marshalled incorrectly!"))
			}

			//ii. Owner's invitation struct changed (updated file uuid) -> save/update Owner's inv struct
			encrypted_owner_struct, encrpyt_and_HMAC_owner_inv_struct, err := EncryptInvitationStruct(userdata, inv_struct, file_flow.Uuid_to_invitation_struct)
			if err != nil {
				userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("Did not encrypt-then_HMAC invitation struct correctly!")))
				return errors.New(strings.ToTitle("Did not encrypt-then_HMAC correctly!"))
			}

			auth_enc_owner_struct := &Encrypted_Struct_or_Contents{
				encrypted_owner_struct,
				encrpyt_and_HMAC_owner_inv_struct,
			}

			//Saving to DataStore
			auth_enc_owner_inv_struct_bytes, e := userlib.Marshal(auth_enc_owner_struct)
			userlib.DatastoreSet(file_flow.Uuid_to_invitation_struct, auth_enc_owner_inv_struct_bytes)
			if e != nil {
				return errors.New(strings.ToTitle("Paniiiccccc"))
			}

			return nil
		}
	}
}
