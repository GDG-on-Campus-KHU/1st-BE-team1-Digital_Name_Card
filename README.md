# Digital Name Card
OAuth2.0을 이용하여 SNS 및 EMail 서비스에 로그인하고 정보를 추출하고

추출한 정보를 JWT로 관리하여 온라인 명함을 생성한다.

자신의 다양한 소셜 계정을 하나의 페이지에서 표현할 수 있다.   
### 실행화면 예시
<img src="https://github.com/GDG-on-Campus-KHU/1st-BE-team1-Digital_Name_Card/blob/master/example.png"  width="450" height="300"/>

## Project Architecture
### Architecture 개요
![프로젝트 구조](https://github.com/GDG-on-Campus-KHU/1st-BE-team1-Digital_Name_Card/blob/master/projectArchitecture.png)

### Flow Chart
<img src="https://github.com/GDG-on-Campus-KHU/1st-BE-team1-Digital_Name_Card/blob/master/SideProject%20Flow%20Chart.png" width="400" height="700"/>

## Requirement
- `go 1.X.X` or higher

## 설치 방법

1. clone source code
    ```
    git clone https://github.com/SeoPPak/GDGoC-BE_SP1.git
    ```
3. move to working directory
    ```
    cd 1st-BE-team1-Digital_Name_Card
    ```
4. install dependencies
    ```
    go mod tidy
    ```

## 실행 방법

1. run Source Code
    ```
    go run .
    ```
    **or**
    ```
    go build -trimpath -ldflags "-w -s" -o server.exe
    .\server.exe
    ```
    
3. Connect to Endpoint
   
    > localhost endpoint
    >
    > <http://localhost:8080/static>

## Configuragion
You must set Environment Variable
Following Environments are needed   

**For Google OAuth**   
`GOOGLE_CLIENT_ID`   
`GOOGLE_CLIENT_SECRET`   
> **Note:** you can check it at your **GCP Console**   

**For Facebook OAuth**   
`FACEBOOK_CLIENT_ID`   
`FACEBOOK_CLIENT_SECRET`   
`FACEBOOK_REDIRECT_URI`   
> **Note:** you can check it at your **Meta for Developer Console**   

You can configure your Environment in your **.env** file 

### Example for .env
```
GOOGLE_CLIENT_ID : "<your_google_client_ID>"
GOOGLE_CLIENT_SECRET : "<your_google_client_secret>"

FACEBOOK_CLIENT_ID=<your_facebook_client_id>
FACEBOOK_CLIENT_SECRET=<your_facebook_client_secret>
FACEBOOK_REDIRECT_URI=http://localhost:8080/auth/facebook/callback
```
   
    
