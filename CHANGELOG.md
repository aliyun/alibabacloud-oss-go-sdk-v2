# ChangeLog - Alibaba Cloud OSS SDK for Go v2

## 版本号：v1.2.3 日期：2025-06-13
### 变更内容
- Fix：Uploader upload from stream fail in special scenarios

## 版本号：v1.2.2 日期：2025-04-25
### 变更内容
- Update：meta query api supports more search condition settings
- Update：uploader supports sequential parameter
- Update：ReadOnlyFile supports OutOfOrderReadThreshold option

## 版本号：v1.2.1 日期：2025-03-07
### 变更内容
- Feature：Add list cloudbox api
- Feature：Supports cloudbox
- Update：Read all the response data before retrying for list api, e.g. ListObjectsV2

## 版本号：v1.2.0 日期：2025-01-08
### 变更内容
- Feature：Add redundancy transition api
- Feature：Add describe regions api
- Break Change：Rename SSEKMSKeyId to ServerSideEncryptionKeyId
- Update：Check the file size before calling truncate
- Fix：Returns all signed header when additional headers is set

## 版本号：v1.1.3 日期：2024-11-29
### 变更内容
- Feature：Add clean restored object api
- Feature：Add add access point for object process
- Update：Add transition time filed for some api.
  
## 版本号：v1.1.2 日期：2024-10-25
### 变更内容
- Feature：Add bucket meta query api
- Feature：Add access point api
- Feature：Add access point public access block api
 
## 版本号：v1.1.1 日期：2024-09-26
### 变更内容
- Fix：Adjust range count when resuming from the last read offset.
- Feature：Add bucket cname api

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