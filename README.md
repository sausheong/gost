# Gost 

[![Go Reference](https://pkg.go.dev/badge/github.com/sausheong/gost.svg)](https://pkg.go.dev/github.com/sausheong/gost)

Gost is short for Go Store. It is a native Go data store for storing data in S3 compatible object storage services. 

In Go we use a lot of structs when we want to do stuff with data. And then when we need to store the data into something more permanent, we break it down into relational database tables or into JSON to store it. When we need it again, we load it from wherever it's stored  back into structs. The data structures we use in our applications are deconstructed and rebuilt each time we save or retrieve from the database.

What if we skip the extra marshalling and unmarshalling step? Just take the struct or string or whichever data type and save it. Then when you need the struct again, just get it back as a struct. No more messing around with marshalling and unmarshalling, converting between different data representations or data types. No more struct tags. Everything's in Go. That's what Gost is about.

<p align="center">
  <img width="200" src="MovingGopher.png" alt="Gopher image from https://github.com/ashleymcnamara/gophers">
  </br>
  <small><a href="https://github.com/ashleymcnamara/gophers">Gopher image from Ashley McNanamara</a></small> 
</p>

Gost is ideal to be used for storing user data. This could be preferences, lists of data the user owns, or anything at all. When a user logs in, he or she can only view his or her own data in a file. Within the file, the keys can be used for different purposes. It could be a map of different kinds of data. The values can be single pieces of data like a string or an int, or it can be a list of data (or a multi-dimensional list of data). It can even be a map. And of course it can any other data structure as well. For example, if you are using a tree in your Go application, you can simply store the entire tree into Gost and get it back as a tree! There is no need to deconstruct and rebuild a tree from a relational database. What more, each user can store different types of data as well, Gost doesn't force any sort of structure at all.

Also, because Gost can store byte arrays it can store images, video, and all sorts of documents as well. Gost can also store documents separately from the data gob and they can be 'published' for public consumption through the cloud storage. For example, you can take and publish images that are publicly available. Gost can also do backups and restores on an individual file basis.

Let's take a closer look at how to use Gost. 


## Using Gost

This section assumes that you already know how to use cloud object stores like Amazon S3, Google Cloud Storage, DigitalOcean Spaces etc. You also know the concepts of access and secret keys as credentials, and the concepts of using storage files into buckets in the object stores.

### Creating a store

Everything in Gost centers around a `Store`. You will use a store to everything else in Gost, so the first thing to do is to create one.

````go
store, err := NewStore(key, secret, endpoint, useSSL, "". bucket)
if err != nil {
    // resolve error
}
````

Let's take a look at the parameters for initialising a `Store` . The `key` is the access key in any one of the object cloud storages. Similarly the `secret` is the secrey key. They will typically come in a pair and you will need to generate them as they are used as credentials to access the cloud storage.

The `endpoint` is the URL used to access the cloud storage (or local storage) and `useSSL` is a boolean that indicates if it uses `http` or `https`. The endpoint for Amazon S3 for example is `s3.amazonaws.com` while for Google Cloud Storage it's `storage.googleapis.com`. For DigitalOcean Spaces it's a bit different, they allow you to set up the location upfront and it's in the endpoint itself. For example, in Singapore I use the `sgp1.digitaloceanspaces.com` endpoint.

The next parameter is the region for which the cloud storage should be hosted. If this is not provided (i.e. an empty string), the default is `us-east-1`. This follows the S3 convention (find the other regions here -- https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.RegionsAndAvailabilityZones.html). In most cases this will match, provided your provider has the same region, though there could be some minor differences. If you're using Google Cloud Storage, you should check out here -- https://cloud.google.com/storage/docs/locations. Depending on where you are or where your server is, you should try to use the nearest provider, for optimal performance.

The last parameter in creating a new store is the `bucket` which is the bucket you want to use to store the data. You can create it in the console or CLI of the cloud storage service you're using, or if you didn't and you specify it here, Gost will create it for you.

Note that if you are using Amazon S3, if you delete your bucket, you have to wait a while before you can create a bucket with the same name.

### Putting data

With the `store` initialized, we can start putting data in. Here's a simple example.

````go
err := store.Put(ctx, "sausheong", "123", "hello world!")
````

The first parameter is the context. You should be using the context for signalling cancellation, timeout, deadlines etc. The parameter `sausheong` is a unique ID, for example it could be the email for a user. The parameter `123` the key, used to reference the data, which is `hello world!`. 

You can put anything in Go, except for channels and functions. Let's take a look at a struct.

````go
thingy := Thingy{
    Name:        "Bob",
    Age:         42,
    DateCreated: time.Now(),
    Length:      1.234,
    Bunch: []OtherThingy{
        {
            Name:   "Alice",
            Number: 1,
        },
        {
            Name:   "Bob",
            Number: 2,
        },
    },
}
Register(thingy)
err := store.Put(ctx, "sausheong", "Bob", thingy)
````

As you see in the code above, you can actually just stuff a nested struct into Gost and it's ok. No more mucking around with JSON, it's all native Go! However before you stuff any custom structs into Gost, you need to register them first so Gost knows what it is. 

````go
Register(thingy)
````

You can register the actual data you want to store, but Gost only really wants to know the struct, so you can just do this.

````go
Register(Thingy{})
````


### Getting data

Getting back the data you stored is quite simple. You just need to know where you stored it, and the key for the piece of data you stored it in.

````go
thing, err := store.Get(ctx, "sausheong", "123")
````

One of the downside of Gost is that you need to know what you stored in there because you need to assert it back to the data type you originally used. In the example above you need to assert `thing` back to `string` because `hello world!` was a string.

````go
thing.(string)
````

You can also get everything back. What you we get is a `map[string]any` and as before you will need to assert it back to whatever it was originally. Remember, if you have stored a custom struct in the store but haven't registered your custom struct, you have to do it before calling `GetAll` because Gost wouldn't know what to do with it.


````go
all, err := store.GetAll(ctx, "sausheong")
thing := all["Bob"].(Thingy) // the custom struct
hello := all["123"].(string) // "hello world!"
````


### Deleting data


Deleting data is quite straightforward, as you might have expected.

````go
err = store.Delete(ctx, "sausheong", "123")
````

You can also delete all the data for a given ID. 

````go
err = store.DeleteAll(ctx, "sausheong")
````

You might wonder why Gost doesn't have anything for updating the data. It's not really necessary because you simply write something else with the same key.

### Storing and retrieving binary data

You might be wondering if Gost can be used to store images or documents like PDF or Microsoft Word files. This is quite trivial for Gost because everything's stored as binary anyway. If you have a document, just open it with Go and make it a byte array, then store the byte array.

````go
imageBytes, err := os.ReadFile("test.png")
err = store.Put(ctx, "sausheong", "test.png", imageBytes)
````

You can do this with any file actually, not just images because in then end it's all byte arrays anyway.

Getting back the file is trivial as well, it's just the reverse of what we just did.

````go
imageBytes, err := store.Get(ctx, "sausheong", "test.png")
err = os.WriteFile("test2.png", imageBytes.([]byte), 0644)
````

Something to note though, all the data is stored in a single file under the same unique ID. If you are planning to store large files, don't store all of them in the same place. Store them under different IDs. Otherwise it's going to be slow everything down.

## Objects

Data in Gost are tied to a unique ID and a key. This allows Gost to provide user-specific data storage where all data related to a specific user (identified by a unique ID) to be stored in a single map. However there are often cases where we need to store data that is common for all users. For example, if we have leaderboard where all users are able to access. This is where objects come in. 

Objects are Gost data that is not identified by a unique ID and a key but by a unique ID only. This single ID is an identifier for the object, instead of being an identifier for a user.

Let's see how this works. Let's say we have a leaderboard, which is just another name for a map of strings to a slice of strings. 

````go
type Leaderboard map[string][]string
board := Leaderboard{}
board["Mona Lisa"] = []string{"Alice", "Bob", "Carol"}
board["The Scream"] = []string{"Dave"}
board["The Starry Night"] = []string{"Carol", "Eve"}
````

Our leaderboard is a ranking of which famous paintings are most liked. We want to allow people to vote for their favorite painting, and also for everyone to view the rankings. 

````go
Register(Leaderboard{})
err = store.PutObject(ctx, "leaderboard", board)
````

As before, we need to register the `Leaderboard` type and simply put the board into Gost, identified by the string "leaderboard".

Getting it out is equally simple.

````go
obj, err := store.GetObject(ctx, "leaderboard")
leaderboard := obj.(Leaderboard)
````

Just get the object back and the assert in back to the `Leaderboard` type.

You can also delete the leaderboard object.

````go
store.DeleteObject(ctx, "leaderboard")
````

**!!IMPORTANT!!** The object functions here are not concurrency-safe. You will likely need to add a mutex or use some other techniques to ensure that race conditions don't appear.

## Publishing files to the Internet

Sometimes you don't want to just store data, you also want the data to be published on the Internet. This is most often used for image files but is also applicable for other types of files like video, PDF, and other documents you want to be directly available. You could of course serve it out from your web application, but why do that when you can use a cloud storage service with a CDN?

### Publishing and unpublishing files

Publishing files are pretty easy. As usual you need to read the file into a byte array and then use the `Publish` function to publish it. You also need to provide a name for the file but more importantly then content type.

````go
imageBytes, err := os.ReadFile("test.png")
loc, err := store.Publish(ctx, "test.png", "image/png", imageBytes)
````
The `Publish` function returns the URL location of the published file. Doing this doesn't automatically make it appear on the Internet though. You need to allow it to be published, which we'll see in the next section. In the meantime let's look at at how a file can be unpublished.

````go
err = store.Unpublish(ctx, "test.png")
````

It's that simple. All published files are put in the `public/` directory as opposed to the `data/` directory for the other data. Also, published files are not identified by a unique ID. 


### Allowing or denying published files to be publicly accessible

As mentioned before once a file is published it's available in the `public/` directory. However this doesn't mean it's accessible on the Internet. To do that you need to allow the `public/` directory to be publicly accessible. You should understand that once the `public/` directory is publicly accessible, all files in it (ie all published files) are as well.

Allowing the `public/` directory to be publicly accessible is simply calling the `AllowPublic` function.

````go
err = store.AllowPublic(ctx)
````

To stop the `public/` directory from being publicly accessible, just call the `DenyPublic` function.

````go
err = store.DenyPublic(ctx)
````

Finally, if you're not sure if it's accessible or not, can just do a quick check. The returned `isPublic` is a boolean that indicates if the `public/` directory is publicly accessible or not.

````go
isPublic, err = store.IsPublic(ctx)
````

## Backing up and restoring

Gost data is always overwritten. To keep a previous copy of the data, you can back it up using the `Backup` function.

````go
err = store.Backup(ctx, "sausheong")
````

This will back up the data identified by the unique ID `sausheong`. There is only 1 backup at any one point in time, so if you call it more than once, it will be overidden.

You can also load the backup and check if there are differences, using the `Load` function.

````go
data, err := store.Load(ctx, "sausheong")
````

Finally you can use the `Restore` function to restore the current data with the backup data.

````go
err = store.Restore(ctx, "sausheong")
````



## Versioning

// TODO

## Encryption

// TODO

## Cloud Storage Services

Here are the specific configurations and gotchas for the different cloud storage services.

### DigitalOcean Spaces

// TODO

### Amazon S3

// TODO

### Google Cloud Storage

// TODO

### Azure Blob Storage

// TODO

### Linode Cloud Storage

// TODO

## Some other tips on using Gost

Gost can be powerful because you can use cloud storage services as your data store and you can access them easily. It's also great to store files that are going to be published to the Internet, because the files are going to be served through the cloud storage services, not by you. The benefits of that are huge -- you can now use very cheap storage that can be managed and secured by someone else. You can also distribute storage not only within the cloud storage provider by using data center regions, but also between multiple cloud storage providers. This makes Gost truly distributed.

Gost is great for storing smaller pieces of information for a specific user when he or she logs in. When the data becomes too large, Gost is inefficient because it needs to load up all the data in memory for use. Also if you are saving the data often, this means a large amount of data could be traveling back and forth from the cloud storage service, which can be slow and definitely not a good thing. With smaller pieces of data this is a lot easier. This doesn't mean you can't use Gost for bigger sets of data, you just need to split it up properly. 

Gost is a lot more useful for non-tabular data where you have Go representations of the data in structs. This is because you can put and get the structs directly! In fact you don't even need to use structs if you can structure your data with the basic maps and slices.

Gost is pretty new, I extracted it from a project I was working on, into a library because I thought it was interesting enough for it to stand alone. This means it's still a work in progress and you shouldn't be surprised if some things don't work out the way it's supposed to.

I tested it mostly with DigitalOcean Spaces, because it has the simplest and easiest console, and also because my other project was hosted there as well. I've tested it on Amazon S3, it works nicely and Google Cloud Storage as well. However Google Cloud Storage doesn't work too well with publishing at the moment, you will need to manually (on the console) set the bucket for public if you want to use it that way.

Gost is not encrypted (yet -- that's a todo). It's serialised in binary but the data is easily exposed if someone gets hold of it. However it's also managed through a cloud storage service so unless you accidentally expose it, you shouldn't need to worry about it.

Performance of Gost improves the nearer it is to the region (obviously). The region setting is important, don't forget that.

Each Gost data file shouldn't be more than 100MB in general, at tops. I tested putting a 100MB data into an empty map, it took 5.8s. Putting another 100MB of data into a map that already has 100MB took 24.3s. Subsequent smaller puts took 17 - 18s. Your mileage might differ. I was interfacing with my laptop running on home broadband going to an S3 endpoint in Singapore. If you do the same from server to server in the same region, things might be quite different. However at the end of the day, lobbing large files back and forth isn't a great idea, so you should keep your files relatively small for best performance.


## Install Minio

Minio is an S3-compatible object storage you can install on your computer. You can install Minio locally and use it for your development purposes. However you should remember that whatever works on your machine doesn't necessarily works the same way on the cloud storage services, so do test the results properly.

### Windows

Install the server by downloading and installing the file from here -- https://dl.min.io/server/minio/release/windows-amd64/minio.exe. Start the server with this command:

````
C:\> .\minio.exe server C:\minio --console-address :9001
````

Install the client by downloading and installing the file from here -- https://dl.min.io/client/mc/release/windows-amd64/mc.exe. You can double click the client to open it up and use it.


### Linux

Install the Minio server.

````
$ wget https://dl.min.io/server/minio/release/linux-amd64/minio
$ chmod +x minio
$ sudo mv minio /usr/local/bin/
````

Install the Minio client.

````
$ wget https://dl.min.io/client/mc/release/linux-amd64/mc
$ chmod +x mc
$ sudo mv mc /usr/local/bin/mc
````

### MacOS

Install the Minio server.

````
$ brew install minio/stable/minio
````

Install the Minio client.  

````
brew install minio/stable/mc
````

## Set up Minio

Start Minio at directory `./data`.

````
$ MINIO_ROOT_USER=admin MINIO_ROOT_PASSWORD=password minio server ./data --console-address ":9001"
````

Set the alias to `local`.

````
$ mc alias set local http://<local host IP>:9000 admin password
````

Set up user service account. Can use any user name.

````
$ mc admin user svcacct add local <username>
````

You can use the key and the secret in Gost.

````
Access Key: <some access key>
Secret Key: <some secret key>
````

To see what the access keys are for this user.

````
mc admin user svcacct list local <username>
````


