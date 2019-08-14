# Lambda image resizer

## What does this do
Takes an s3 object (let's use the term folder) and resizes and formats all the images it contains
- It only processes one level e.g. images in `/quidco-images/banner/` will be processed, `/quidco-images/banner/mobile` 
won't be
- Creates a folder in the same directory for the resized images e.g. `/quidco-images/banner/1100x250`
- It will only resize an image to a specified dimension if the image is larger than the dimension 
e.g. if the specified dimensions are `1100x250` but the image size is `1000x100`, that image will be skipped
- It will crop an image to fit the dimensions (instead of squeeze the image and make it look distorted) 
e.g. using the the dimensions `1100x250` as our target, it will resize a `1920x1080` image to a `1100x619` image and trim 
the the excess height evenly from both top and bottom to give us a final size of `1100x250`

## Running the application

### Prerequisites
- Install golang https://golang.org/doc/install#install
- Install docker https://docs.docker.com/docker-for-mac/install/
- Pull the repo using `go get github.com/mkadiri/lambda` (recommended as this will place the project in the correct folder)

You should now find the project in `~/go/src/github.com/mkadiri/lambda`

### Use the Makefile 

In the `Makefile` you'll find two commands `make run` and `make zip`

#### make run 

The docker-compose has everything you need to run this application locally, however **you'll need to specify your AWS credentials**.

For testing purposes you'll need to create an s3 bucket with a folder, with images

The function name `app` as well as an event json will need to be passed.

The event json needs to specify the s3 bucket, s3 object (folder) and also the width and height you want to resize images to.

**Note: make sure `folder` ends with a trailing slash**

 
 ```
 {
    "bucket": "quidco-images",
    "folder": "merchant-cover-photos/banners/",
    "width": 1100,
    "height": 250
  }
```

#### make zip

The application can be built and zipped with the `make zip` command, this zip file can be uploaded in to AWS
      
## Resources
- s3 bucket: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/s3-example-basic-bucket-operations.html
