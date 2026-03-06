# AWS S3 STS AssumeRole 配置说明

## 背景

S3 预签名 URL 默认会在 `X-Amz-Credential` 参数中携带 Access Key ID（AK）。
通过配置 STS AssumeRole，可将永久凭证（`AKIA` 开头）替换为临时凭证（`ASIA` 开头），临时凭证到期后自动失效，降低泄露风险。

## AWS 侧配置（一次性操作）

### 1. 创建 IAM Role

在 AWS 控制台 → IAM → 角色 → 创建角色，选择**自定义信任策略**：

```json
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Principal": {
      "AWS": "arn:aws:iam::<账号ID>:user/<IAM用户名>"
    },
    "Action": "sts:AssumeRole"
  }]
}
```

### 2. 给 Role 附加 S3 权限策略

```json
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": ["s3:PutObject", "s3:GetObject", "s3:DeleteObject", "s3:ListBucket"],
    "Resource": [
      "arn:aws:s3:::<bucket名称>",
      "arn:aws:s3:::<bucket名称>/*"
    ]
  }]
}
```

### 3. 给 IAM User 添加 AssumeRole 权限

```json
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": "sts:AssumeRole",
    "Resource": "arn:aws:iam::<账号ID>:role/<角色名称>"
  }]
}
```

---

## 应用侧配置

### 方式一：环境变量（推荐）

```bash
export OssRoleArn=arn:aws:iam::350646758444:role/tier0-oss-s3-role
```

Docker / K8s 示例：
```yaml
env:
  - name: OssRoleArn
    value: arn:aws:iam::350646758444:role/tier0-oss-s3-role
```

### 方式二：yaml 配置文件

```yaml
OssConf:
  OssType: aws
  AccessKeyID: AKIAVDJBM5QWLQ67DDSL
  AccessKeySecret: <SecretKey>
  Region: ap-southeast-1
  Location: s3.amazonaws.com
  TemporaryBucketName: tier0-upload-temp-pre
  PublicBucketName: tier0-upload-pub
  PrivateBucketName: tier0-upload-private
  RoleArn: arn:aws:iam::350646758444:role/tier0-oss-s3-role
  # RoleSessionName: oss-session   # 可选，默认 oss-session
  # SessionDuration: 3600          # 可选，单位秒，范围 900~43200
```

> 环境变量优先级高于配置文件。

---

## 效果对比

| 配置 | URL 中的 AK | 是否会过期 |
|------|------------|----------|
| 未配置 RoleArn | `AKIA...`（永久凭证） | 否 |
| 配置 RoleArn | `ASIA...`（临时凭证） | 是，按 SessionDuration 过期 |

---

## 运行测试

```bash
# 静态凭证测试
AWS_ACCESS_KEY_ID=xxx AWS_SECRET_ACCESS_KEY=xxx \
S3_ENDPOINT=s3.amazonaws.com S3_REGION=ap-southeast-1 S3_BUCKET=your-bucket \
go test ./oss/... -run "TestSignedPutUrl_StaticCreds|TestSignedGetUrl_StaticCreds" -v

# STS 临时凭证测试
AWS_ACCESS_KEY_ID=xxx AWS_SECRET_ACCESS_KEY=xxx \
S3_ENDPOINT=s3.amazonaws.com S3_REGION=ap-southeast-1 S3_BUCKET=your-bucket \
OssRoleArn=arn:aws:iam::xxx:role/xxx \
go test ./oss/... -run "TestSignedPutUrl_STS" -v
```
