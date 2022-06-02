# SAD Project API Documentation

[TOC]

## **API Base URL** 
https://sadbackend-cyt.herokuapp.com
## **API Testing Interface** 
https://sadbackend-cyt.herokuapp.com/swagger/index.html#/

## 跟後端要 available time
> [name=cytttt] **進入這個頁面之前你應該要打 `GET /api/v1/gym/list ` 決定你可以顯示哪些分店供用戶選擇**
> [name=cytttt] **請看清楚 `GET /api/v1/gym/list ` 中回傳的 `available_machine` 的內容**
> [name=cytttt] **根據 user 選擇的 part 你要在 dropdown 顯示他可以選的 machine**
> [name=cytttt] **然後你要在 dropdown 顯示 `"name"` 然後你要記錄 `"branch_gym_id"`**
> [name=cytttt] **現在只有 branch_gym_id=branch-1000001, mahcine = `"treadmill"`(有氧) or `"pec deck"`(練胸)**
- **Endpoint**:  `GET /api/v1/user/available`
- **Request**: query string 有 5 個欄位需要填值，請注意下方 url ? 後面的內容
    **請試用 swagger 測試介面**
    - date
    - period 
    - account
    - branch_gym_id
    - machine
    - **request url example**
`https://sadbackend-cyt.herokuapp.com/api/v1/user/available?date=2022-06-29&period=morning&account=meowmeow123&branch_gym_id=branch-1000001&machine=treadmill`
- **Response**:
```json=
// unsuccessful lookup
// data 是 null 就代表找不到可以的時間
{
  "code": 200,
  "msg": "Ok",
  "data": null
}
// successful lookup
{
  "code": 200,
  "msg": "Ok",
  "data": [ 
    {
      "start": "06:30",
      "end": "07:00",
      "machine_id": "machine-1000006"
    },
    {
      "start": "07:00",
      "end": "07:30",
      "machine_id": "machine-1000006"
    },
    {
      "start": "08:00",
      "end": "08:30",
      "machine_id": "machine-1000006"
    },
    {
      "start": "09:00",
      "end": "09:30",
      "machine_id": "machine-1000006"
    },
    {
      "start": "10:00",
      "end": "10:30",
      "machine_id": "machine-1000005"
    }
  ]
}

Response headers
```
training part 部位: "back", "chest", "cardio", "abs", "leg", "arm", "hips"

## make reservation
> [name=cytttt] 你應該要記錄 available time 回傳的 machine_id 取得 且你只應該使用 available time 回傳的 machine_id 因為並不是所有 machine 都可以預約
> [name=cytttt] machine-1000005, machine-1000006, machine-1000007 這三個是可以預約的
- **Endpoint**:  `POST /api/v1/user/reservation`
- **Request**: example
```json=
{
  "account": "meowmeow123",        // 用戶帳號
  "date": "2006-01-02",            // 2006/01/02
  "machine_id": "machine-1000007", // machine id 
  "start": "13:00"                 // 預約起始時間
}
```
- **Response**:
```json=
{
  "code": 200,
  "msg": "Ok",
  "data": nil,
}
```

## 獲得各分店的人潮狀況
> [name=cytttt] **請看清楚是"[ ]"還是"{ }"**
- **Description**: staff 的 gym list 頁面，拿當下各分店的人潮狀態
- **Endpoint**: `GET /api/v1/gym/list `
- **Request**: no
- **Response**: 
```json=
{
  "code": 200,
  "msg": "Ok",
  "data": [
      {
          "branch_gym_id": "branch-1000001",
          "name": "Daan",
          "address": "No. 130-2, Sec. 1, Fuxing S. Rd., Da’an Dist., Taipei City",
          "status": "crowded",
          "available_machine": {
            "abs": [
                "ab bench",
                "ab coaster"
            ],
            "arm": [
              "bicep vest",
              "rower"
            ],
            "back": [
              "pulldown",
              "power tower"
            ],
            "cardio": [
              "spinner",
              "treadmill"
            ],
            "chest": [
              "chest press",
              "pec deck"
            ],
            "hips": [
              "donkey kick",
              "ham raise"
            ],
            "leg": [
              "hack squat",
              "leg press"
            ]
          }
        }
      ]
    }
  ]
}

```

## 登記器材人數
> [name=cytttt] crowded 的部分建議打 api/v1/gym/list
- **Description**: staff 登記個器材人數頁面，按下 + 或 - 的按鈕後，會回傳更新後的各分店器材的人數資訊(一次只會改一個機器的人數，但是要回傳所有分店機器的資料)
- **Endpoint**: `PUT /api/v1/machine/status`
- **Request**:
```json=
{
  "amount": 1,
  "machine_id": "machine-1000001"
}
```
- **Response**:
```json=
{
  "code": 200,
  "msg": "Ok",
  "data": [
    {
      "machine_id": "machine-1000001",
      "name": "treadmill-101",
      "waiting_ppl": 4,
      "category": "cardio"
    },
    {
      "machine_id": "machine-1000002",
      "name": "treadmill-102",
      "waiting_ppl": 4,
      "category": "cardio"
    },
    {
      "machine_id": "machine-1000003",
      "name": "power tower 1",
      "waiting_ppl": 5,
      "category": "back"
    }
  ]
}
```


## 獲得某家分店所有 machine 的狀態
- **Endpoint**: `GET /api/v1/gym/machine`
- **Resquest**: query string gym_id, sorted_by
- **Response**:
```json=
{
  "code": 200,
  "msg": "Ok",
  "data": [
    {
      "machine_id": "machine-1000001",
      "name": "treadmill-101",
      "waiting_ppl": 4,
      "category": "cardio"
    },
    {
      "machine_id": "machine-1000002",
      "name": "treadmill-102",
      "waiting_ppl": 4,
      "category": "cardio"
    },
    {
      "machine_id": "machine-1000003",
      "name": "power tower-1",
      "waiting_ppl": 5,
      "category": "back"
    }
  ]
}
{
  "code": 200,
  "msg": "Ok",
  "data": {
    "back": [
      {
        "machine_id": "machine-1000003",
        "name": "power tower-1",
        "waiting_ppl": 5,
        "category": "back"
      }
    ],
    "cardio": [
      {
        "machine_id": "machine-1000001",
        "name": "treadmill-101",
        "waiting_ppl": 4,
        "category": "cardio"
      },
      {
        "machine_id": "machine-1000002",
        "name": "treadmill-102",
        "waiting_ppl": 4,
        "category": "cardio"
      }
    ]
  }
}
```
## 註冊
- **Description**: 傳入一組：帳號、密碼、姓名、性別、email、電話、生日、身高、體重，回傳此組帳密是否有人使用；若沒有，回傳一個判斷為正確的訊息；若有則回傳一個判斷為錯誤的訊息。
- **Endpoint**: `POST /api/v1/user/signup`
- **Request**:
```json=
{
  "account": "meowmeow789",
  "birthday": "2006/01/02",
  "email": "meowantony@gmail.com",
  "gender": "male",
  "height": 188.87,
  "name": "Antony Cho",
  "password": "meowmoew22",
  "phone": "0912345678",
  "weight": 69.69
}
```
- **Response**:
```json=
{
  "code": 200,
  "msg": "Ok",
  "data": null
}
{
  "code": 10003,
  "msg": "User already exist",
  "data": null
}
```
## 登入
> [name=cytttt] 我設成同時支援 client, staff 的登入 根據提供的 `userrole`。account 可以傳 account 或 email。不同的 error 會有不同的`code, msg` 請自行確認:) 
- **Description**: 傳入一組帳密，確認這組帳密是否有註冊過；若有，回傳一個判斷為正確的訊息；若沒有則回傳一個判斷為錯誤的訊息(有帳號但沒密碼/沒有帳號為兩種不同的錯誤)
- **Endpoint**: `POST /api/v1/user/login`
- **Request**:
```json=
{
  "account": "meowmeow123",
  "password": "meowmoew22",
  "user_role": "client"
}
```
- **Response**:
```json=
{
  "code": 200,
  "msg": "Ok",
  "data": {
    "account": "meowmeow123",
    "name": "testMeowClient",
    "user_role": "client",
    "level": ""
  }
}
```
- **Testing Accounts**:
```json=
{
  "account": "meowmeow123",
  "password": "meowmoew22",
  "user_role": "client"
}
```
```json=
{
  "account": "meowmeow456",
  "password": "meowmoew22",
  "user_role": "staff"
}
```
## 更改個人資訊
- **Description**: 更改個人資訊，將更改內容回傳至後端，會回傳成功/失敗
- **Endpoint**: `PUT /api/v1/user/info`
- **Resquest**:
```json=
{
  "account": "meowmeow123",
  "day": 29,
  "email": "meowtestclient@gmail.com",
  "gender": "male",
  "height": 180.13,
  "month": 5,
  "name": "testMeowClient",
  "pay_type": "visa",
  "payment_plan": "1234123412341234",
  "phone": "0919886886",
  "plan": "normal",
  "weight": 69.69,
  "year": 2001
}
```
- **Response**:
```json=
{
  "code": 200,
  "msg": "Ok",
  "data": {
    "account": "meowmeow123",
    "name": "testMeowClient",
    "email": "meowtestclient@gmail.com",
    "personal_info": {
      "Gender": "male",
      "Phone": "0919886886",
      "Birthday": "2001-05-28T16:00:00Z"
    },
    "body_info": {
      "Height": 180.42,
      "Weight": 69.69
    },
    "subscription": {
      "Plan": "",
      "ExpiredAt": "0001-01-01T00:00:00Z"
    },
    "payment_method": {
      "PayType": "",
      "Account": ""
    },
    "AttendenceRecord": null,
    "created_at": "0001-01-01T00:00:00Z",
    "updated_at": "2022-05-14T07:38:17.839Z"
  }
}
```
## 獲得個人資訊
- **Endpoint**: `GET /api/v1/user/info`
- **Resquest**: query string
- **Response**:
```json=
{
  "code": 200,
  "msg": "Ok",
  "data": {
    "account": "meowmeow123",
    "name": "testMeowClient",
    "email": "meowtestclient@gmail.com",
    "personal_info": {
      "Gender": "male",
      "Phone": "0919886886",
      "Birthday": "2001-05-28T16:00:00Z"
    },
    "body_info": {
      "Height": 180.42,
      "Weight": 69.69
    },
    "subscription": {
      "Plan": "",
      "ExpiredAt": "0001-01-01T00:00:00Z"
    },
    "payment_method": {
      "PayType": "",
      "Account": ""
    },
    "AttendenceRecord": null,
    "created_at": "0001-01-01T00:00:00Z",
    "updated_at": "2022-05-14T07:38:17.839Z"
  }
}
```
## Staff 統計頁面
> [name=cytttt] ~~我對於 attendence 有疑問想 5/15 釐清一下~~
> 希望能傳當前時區的 query 或帶 header 之類的，理由同下面那隻 api 
> 目前寫死 utc+8 但我不太舒服
- **Description**: staff的 statistic 頁面，拿過去七天(收到 request 當天不算)的用戶出席次數和平均用戶停留時間
- **Endpoint**: `GET /api/v1/user/staff/stat`
- **Resquest**:
- **Response**:
```json=
{
  "code": 200,
  "msg": "Ok",
  "data": [
    {
      "date": "2022/05/11",
      "attendance_count": 2,
      "avg_stay_hour": 1.6111112
    },
    {
      "date": "2022/05/12",
      "attendance_count": 2,
      "avg_stay_hour": 1.5277778
    },
    {
      "date": "2022/05/13",
      "attendance_count": 1,
      "avg_stay_hour": 0.5
    },
    {
      "date": "2022/05/14",
      "attendance_count": 2,
      "avg_stay_hour": 1.7777778
    },
    {
      "date": "2022/05/15",
      "attendance_count": 1,
      "avg_stay_hour": 2
    },
    {
      "date": "2022/05/16",
      "attendance_count": 2,
      "avg_stay_hour": 1.25
    },
    {
      "date": "2022/05/17",
      "attendance_count": 1,
      "avg_stay_hour": 1
    }
  ]
}

Response headers
```
## Home 界面獲取 reservation 資料 
> [name=cytttt] 目前應該沒有 activity, 時間請自行轉換成台灣的時區(我給的會是utc+0) 後端理應不知道開網頁的人在哪個時區 time format: RFC3339
> [name=cytttt] **不會回傳過期的 reservation ，回傳將來最近的前四筆**
- **Description**:User進入頁面得到的資料
- **Endpoint**:  `GET /api/v1/user/reservation/:account`
- **Request** : 
- **Response**: 
```json=
{
  "code": 200,
  "msg": "Ok",
  "data": [
    {
      "category": "cardio",
      "machine_id": "machine-1000001",
      "machine_name": "treadmill-101",
      "gym_id": "branch-1000001",
      "gym_name": "Daan",
      "date": "2023-12-09T16:00:00Z"
    }
  ]
}
```

## 用戶統計頁面 Client Stat 
>[name=kajie] 這個data也是那個{}的問題麻煩你改一下，謝謝。
> [name=cytttt]不是 他們 type (int, string)不一樣根本就放不了進 array 啊??
> [name=kaijie] 問題解決了~
- **Description**:User Gym stat
- **Endpoint**:  `GET /api/v1/user/stat/:account`
- **Request** : 
- **Response**: 
```json=
{
  "code": 200,
  "msg": "Ok",
  "data": {
    "stay_time": 232,
    "calories": 32424,
    "most_train": "cardio",
    "least_train": "arm"
  }
}
```

## 客戶獲得 gym 內 machine 的資料 categorized by training part 
> [name=cytttt] 可以參考 `GET /api/v1/gym/machine` 我有更新你昨天的說的需求
> [name=kajie]回復:可把格式改成我給你的這個嗎？或是把你的data的'{}'改成‘[]’,這個問題我花了幾個小時才找到。。。，因為不是object就讀取不到。。。，但還是希望你改成我的格式，麻煩你了。還有我想問你為什麼你的data資料格式會跟getgymlist 的不同,會想用'{}'
> [name=cytttt] re: 還有我想問你為...,會想用'{}'? 
> 因為按照部位分類，啊我的理解是這樣被每個部位之下才是一個 array 。 
> 假如用你的 array 的 format 你要怎麼知道 cardio 在第幾個？這樣你是不是要先 traverse 一遍 response['data']
> 然後才知道喔 "cardio" 是 array 的第一個 element 然後才知道用 response['data'][1] 去拿 "cardio" 的 machine list
> 我現在的寫法用 object 的話就沒有這問題 直接用 response['data']['cardio'] 就可以拿 "cardio" 的 machine list。
> [name=cytttt] re: 可把格式改成我給你的這個嗎？ 
> "gym" 是不是可以刪掉？ 請問 "id" 的用途？ "id" 跟 "workout" 的 mapping?
> [name=cytttt] 還是你比較希望用 url path 的方式傳參數而不是用 query parameter？
- **Description**:User Gym Machines
- **Endpoint**:  `GET /api/v1/gym/machine/category/:gym_id`
- **Request** : 
- **Response**: 
```json=
{
  "code": 200,
  "msg": "Ok",
  "data": [
    {
      "category": "back",
      "machines": [
        {
          "machine_id": "machine-1000003",
          "name": "power tower 1",
          "waiting_ppl": 5,
          "category": "back"
        }
      ]
    },
    {
      "category": "chest",
      "machines": null
    },
    {
      "category": "cardio",
      "machines": [
        {
          "machine_id": "machine-1000001",
          "name": "treadmill-101",
          "waiting_ppl": 4,
          "category": "cardio"
        },
        {
          "machine_id": "machine-1000002",
          "name": "treadmill-102",
          "waiting_ppl": 4,
          "category": "cardio"
        }
      ]
    },
    {
      "category": "abs",
      "machines": null
    },
    {
      "category": "leg",
      "machines": null
    },
    {
      "category": "arm",
      "machines": null
    },
    {
      "category": "hips",
      "machines": null
    }
  ]
```