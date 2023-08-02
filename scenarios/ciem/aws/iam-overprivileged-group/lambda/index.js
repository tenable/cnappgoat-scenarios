const AWS = require('aws-sdk');

exports.handler = async function(event, context) {
    const secretsmanager = new AWS.SecretsManager();
    const secretId = event.SecretId;
    const bucketName = event.BucketName;
    const accessKeyId = event.AccessKeyId;
    console.log(`SecretId: ${secretId}`);
    console.log(`BucketName: ${bucketName}`);
    console.log(`AccessKeyId: ${accessKeyId}`);

    try {
        const data = await secretsmanager.getSecretValue({ SecretId: secretId }).promise();
        console.log(data);
        const accessKeySecret = data.SecretString;

        const s3 = new AWS.S3({
            accessKeyId: accessKeyId,
            secretAccessKey: accessKeySecret,
        });

        const bucketContents = await s3.listObjectsV2({
            Bucket: bucketName,
        }).promise();

        console.log(bucketContents);
    } catch (err) {
        console.error(err);
    }
}
