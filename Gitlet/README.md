# Gitlet Design Document

**By**: Iman Kahssay

## Classes and Data Structures

### Class 1: Main

This is the entry point to our program.
It takes in arguments from the command line and based on the 
command (the first element of the args array) calls the 
corresponding command in Repository which will 
actually execute the logic of the command. 
It also validates the arguments based on the command to 
ensure that enough arguments were passed in.

#### Fields

This class has no fields and hence no associated state: it 
simply validates arguments and defers the execution to the 
CapersRepository class.


### Class 2: Repository

This is where the main logic of our program will live. This 
class will handle all of the actual Gitlet commands by 
reading/writing from/to the correct file, setting up 
persistence, and additional error checking.

It will also be responsible for setting up all persistence 
within gitlet. This includes creating the .gitlet folder as 
well as the folder and file where we store all Commit objects 
and the Blobs.

This class defers all Commit specific logic to the Commit class: 
for example, instead of having the Repository class handle 
Commit serialization and deserialization, we have the Commit 
class contain the logic for that.


#### Fields

1. Field 1: **static final File CWD = new File(System.getProperty("user.dir"))**
   
    The Current Working Directory. Since it has the 
    package-private access modifier (i.e. no access modifier), 
    other classes in the package may use this field. 
    It is useful for the other File objects we need to use.

2. Field 2: **static final File GITLET_REPOSITORY_FOLDER = Utils.join(CWD, ".gitlet")** 
   
    The hidden .gitlet directory. This is where all of the 
    states of the Repository will be stored, including additional 
    things like the Commit objects and the Blobs. It is also 
    package private as other classes will use it to store their state.

3. Field 3: **public File HEAD = join(Repository.GITLET_REPOSITORY, "HEAD")**

   'Points' to the most current commit (the 'leaf') by saving the
    SHA-1 Hash of that Commit (i.e. commitID) in a text file.
   
4. Field 4: **public static final File MASTER = join(Repository.GITLET_REPOSITORY, "MASTER")**

   'Points' to the 'root' commit by saving the SHA-1 Hash of 
   that commit (i.e. commitID) (therefore, saving the branch/tree 
   of the repository).

These fields are both static since we don’t actually instantiate 
a Repository object: we simply use it to house functions. If 
we had additional non-static state (like the Dog class), we’d 
need to serialize it and save it to a file.


### Class 3: Commit

This class represents a Commit that will be stored in a file. 
Because each Commit will have a unique ID, we may simply use 
that as the name of the file that the object is serialized to.

All Commit objects are serialized within the COMMITS_FOLDER which 
is within the GITLET_REPOSITORY_FOLDER. The Commit class has 
helpful methods that will return the Commit object corresponding 
to some String ID given to it, as well as write that Commit to 
a file to persist its changes.

#### Fields

1. Field 1: **static final File COMMITS_FOLDER = Utils.join(Repository.GITLET_REPOSITORY_FOLDER, ".commits")** 
   
   The File object that corresponds to the directory containing 
   all the serialized Commit objects. This is static since all 
   Commit objects are stored within the same directory. 

2. Field 2: **private long timestamp** -> The time Commit object was "committed" (or created).

3. Field 3: **private String message** -> User's commit message when creating a Commit.

4. Field 4: **private String commitID** -> The SHA1 hashvalue of Commit(will also be the File Object's name strong this Commit Object).

5. Field 5: **private String commitParentID** -> The SHA1 hashvalue of previous Commit (i.e. commit's parent).

6. Field 6: **private Hashtable<String, String> tableOfBlobs;** ->  Hashtable of Blob references (the SHA1 value of Blobs (it's ID)).

### Class 4: Blob

This class represents a Blob that will be stored in a file.
Because each blob will have a unique ID, we may simply use that
as the name of the file that the object is serialized to.
This class will save the different *versions* of a single file.

#### Fields

1. Field 1 : **static final File BLOBS_FOLDER = Utils.join(Repository.GITLET_REPOSITORY, ".blobs")**
   
    The File object that corresponds to the directory 
    containing all the serialized Blob objects.
    This is static since all Blob objects are stored within 
    the same directory.
   
2. Field 2: **private String blobID** -> The SHA1 value of this Blob (it's ID).

3. Field 3: **private String fileName** -> Name of the file Blob will be saving the contents of.
  
4. Field 4: **private File fileContents** -> Saves the target files version/contents into a text file.
    - Thought about saving it as a String, but the contents may be too large.

### Class 5: Utils

Utils
This class contains helpful utility methods to read/write 
objects or String contents from/to files, as well as reporting 
errors when they occur.

This is a staff-provided and PNH written class, so we leave 
the actual implementation as magic and simply read the helpful 
javadoc comments above each method to give us an idea of 
whether or not it’ll be useful for us.

### Class 6: Staging Area

This class represents the Stating Area that stores files before
they are committed (in add and remove commands).

#### Fields

1. Field 1 :

### Class 7: Staging Area

This class represents the Stating Area that stores files before
they are committed (in add and remove commands).

#### Fields

1. Field 1 :

### Class 8: Gitlet Exception

This class represents all of the exceptions and/or error 
messages that may be thrown.

#### Fields

1. Field 1 :

## Algorithms

## Persistence

The directory structure looks like this:

    CWD                                         <==== Whatever the current working directory is.
    └── .gitlet                                 <==== All persistant data is stored within here
        ├── commits                             <==== All commits are stored in this directory
            ├── commit1_ID                      <==== A single Commit instance stored to a file
            ├── commit2_ID
            ├── ...
            └── commit_n_ID
        └── blobs                               <==== All blobs are stored in this directory
            ├── blob1_ID                        <==== A single Blob instance stored to a file
            ├── fileContents for Blob1_ID       <==== The fileContents saved in a file in a Blob object.
            ├── blob2_ID 
            ├── fileContents for Blob2_ID 
            ├── ...
            ├── blob_n_ID
            └── fileContents for Blob_n_ID 
        └── HEAD.txt                              <==== Persistent Head Pointer saved in a text file (like story.txt in lab 6).
       

