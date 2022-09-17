# Gost - Native Go data store on Cloud Object Storages

Gost is a native Go data store for storing data in S3 compatible object stores. 

## Object storages

Object storage is a popular form of data storage where data is managed as objects as compared with other storage architectures. Object storage is very popularly used for cloud storage. The most popular service today is Amazon Simple Storage Service (S3) and it has more or less become the defacto standard for cloud object storage services. Other object cloud storage services include Google Cloud Storage, Azure Blob Storage, DigitalOcean Spaces, Oracle Cloud Object Storage, Linode Object Storage and so on. There are also open source object storage software that you can install on your own machines including OpenStack Swift, Minio, Zenko CloudServer, OpenIO and many others. Interesting all of them are more or less S3 compatible, which shows the great strength of S3 as a service.

## Gost

Gost is short for Go Storage.

It's an interesting take for storing data for Go applications. In Go we use a lot of structs when we want to do stuff with data. And then when we need to store the data into something more permanent, we break it down into JSON, into relational database tables and so on, and store it. When we need it back, we load it back into structs and use it.

What if we skip the extra marshalling and unmarshalling step? Just take the struct and save it. Go has a mechanism for it -- it's called gob. It's a binary serialization package used to create streams of binary data that can be used in RPCs. It can also be used to scrunch structs and any Go data types into a binary, self-describing form. Gob works for everything except channels and functions. 

In Gost, I take whatever Go data type you have, and I stuff them into a map with a string key and value of any data type (I literally use `any`). Then I take a unique ID formed by doing a base64 of an ID value you pass in (for example, an email address) and that becomes the filename of the gob file. That gob file, which was originally stored in the filesystem, is now in an object store of your choice, under the `data` directory.

Why is it done this way? Gost is ideal to be used for storing user data. This could be preferences, lists of data the user owns, or anything at all. When a user logs in, he or she can only view his or her own data in this file. Within the file, the keys can be used for different purposes. It could be a map of different kinds of data. The values can be single pieces of data like a string or an int, or it can be a list of data (or a multi-dimensional list of data). It can even be a map. And of course it can a list of structs etc as well. What more, each user can store different types of data as well, there isn't really a need for all data to be the same. 

With this, you can imagine that because it can store byte arrays it can store images, video, and all sorts of documents as well. One of the interesting things about using cloud object storage services, Gost can store documents separately from the data gob and they can be 'published' for public consumption. For example, you can take and publish images that are publicly available. We'll take a second look at this later.

Gost is terrible if you want to use it to do analytics on the user data because the data is essentially unstructured. But you normally wouldn't want to do analytics on transaction data anyway, you would want to suck them all up and then throw it into another data store for better analysis.

Let's take a closer look at how we to use Gost.


## Using Gost

### Creating a store

Everything in Gost centers around a `Store`. You will use a store to everything else in Gost, so the first thing to do is to create one.

````go
store, err := NewStore(key, secret, endpoint, useSSL, bucket)
if err != nil {
    // resolve error
}
````

Let's take a look at the few variable. The `key` is the access key in any one of the object cloud storages. Similarly the `secret` is the secrey key. They will typically come in a pair and you will need to generate them as they are used as credentials to access the cloud storage.

The `endpoint` is the URL used to access the cloud storage (or local storage) and `useSSL` is a boolean that indicates if it uses `http` or `https`. The endpoint for Amazon S3 for example is `s3.amazonaws.com` while for Google Cloud Storage it's `storage.googleapis.com`.

The last parameter in creating a new store is the `bucket` which is the bucket you want to use to store the data. You can create it in the console or CLI of the cloud storage service you're using, or if you didn't and you specify it here, Gost will create it for you.

### Putting data

With the store, we can start putting data in. Here's a simple example.

````go
err := store.Put(ctx, "sausheong", "123", "hello world!")
````

The first parameter is the context. You should be using the context for signalling cancellation, timeout, deadlines etc. The parameter `sausheong` is the user ID for the user, for example it could be his email. The parameter `123` the key, used to reference the data, which is `hello world!`. 

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

When you register you can register the actual data you want to store, but Gost only really wants to know the struct, so you can just do this.

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

## When to use and when NOT to use Gost

Gost can be powerful because you can use cloud storage services as your data store and you can access them easily. It's also great to store files that are going to be published to the Internet, because the files are going to be served through the cloud storage services, not by you. The benefits of that are huge -- you can now use very cheap storage that can be managed and secured by someone else. You can also distribute storage not only within the cloud storage provider by using data center regions, but also between multiple cloud storage providers. This makes Gost truly distributed -- 

Gost is great for storing smaller pieces of information for a specific user when he or she logs in. When the data becomes too large, Gost is inefficient because it needs to load up all the data in memory for use. Also if you are saving the data often, this means a large amount of data could be traveling back and forth from the cloud storage service, which can be slow and definitely not a good thing. With smaller pieces of data this is a lot easier. This doesn't mean you can't use Gost for bigger sets of data, you just need to split it up properly. 

Gost is a lot more useful for non-tabular data where you have Go representations of the data in structs. This is because you can put and get the structs directly! In fact you don't even need to use structs if you can structure your data with the basic maps and slices.

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


