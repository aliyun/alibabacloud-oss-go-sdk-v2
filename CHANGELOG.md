# ChangeLog - Alibaba Cloud OSS SDK for Go v2

## 版本号：v1.1.0 日期：2024-09-18
### 变更内容
- Feature：Add bucket logging api
- Feature：Add bucket worm api
- Feature：Add bucket policy api
- Feature：Add bucket transfer acceleration api
- Feature：Add bucket archive direct read api
- Feature：Add bucket website api
- Feature：Add bucket worm api
- Feature：Add bucket https configuration api
- Feature：Add bucket resource group api
- Feature：Add bucket tags api
- Feature：Add bucket encryption api
- Feature：Add bucket referer api
- Feature：Add bucket inventory api
- Feature：Add bucket access monitor api
- Feature：Add bucket style api
- Feature：Add bucket replication api
- Feature：Add bucket cors api
- Feature：Add bucket lifecycle api
- Feature：Add public access block api
- Update：Add more fileds for bucket info api
- Update：Add more fileds for bucket stat api
- Update：Refine filepath check in uploader
- Update：Add a slash when building the path-style request url.
- Fix：Fix FileExists and DirExists function bug
- Break Change：Modify the Tagging.TagSet type from TagSet to *TagSet. 

## 版本号：v1.0.2 日期：2024-06-25
### 变更内容
- Fix：resumable download/upload progress bug
- Fix：resumable upload bug
 
## 版本号：v1.0.1 日期：2024-06-21
### 变更内容
- Feature：Add additional Headers opition
- Feature：Add user agent opition
- Feature：Add the expiration check for presign api
- Document：Added english version of the developer guide
 
## 版本号：v1.0.0 日期：2024-05-29
### 变更内容
- Feature：Add credentials provider
- Feature：Add retryer
- Feature：Add signer v1 and signer v4
- Feature：Add httpclient
- Feature：Add bucket's basic api
- Feature：Add object's api
- Feature：Add presigner
- Feature：Add paginator
- Feature：Add uploader, downloader and copier
- Feature：Add file-like api
- Feature：Add encryption client
- Feature：Add IsObjectExist/IsBucketExist api
- Feature：Add PutObjectFromFile/GetObjectToFile api
- Document：Add developer guide