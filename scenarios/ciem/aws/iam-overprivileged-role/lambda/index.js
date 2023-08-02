const AWS = require('aws-sdk');

exports.handler = async function(event, context) {
    console.log("Received event: ", event);

    // Parse the BucketName from the event
    const bucketName = event.BucketName;

    const s3 = new AWS.S3();
    // List all buckets
    const buckets = await s3.listBuckets().promise();

    // List the bucket
    const data = await s3.listObjectsV2({
        Bucket: bucketName
    }).promise();
    // Get bucket public access block
    const publicAccessBlock = await s3.getPublicAccessBlock({
        Bucket: bucketName
    }).promise();
    console.log("Public access block: ", publicAccessBlock)
    console.log("Buckets: ", buckets);
    console.log("Bucket contents: ", data);
}
