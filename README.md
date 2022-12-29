# Back-End No Stunting App

> This app is using
>
> 1. Golang
> 2. Gin-Gonic
> 3. MongoDB
> 4. Gorilla Websocket

***

## Authentication

1. Super Admin
2. Admin
3. Facility
4. Mother
5. Child

### Super Admin

Super admin has fully control of REST-API

### Admin

Admin can control facility, mother and child who registered under its control.

### Facility

Facility can control mother and child who registered under its control.

### Mother

Mother only be able to control their own monitoring field.

### Child

Child who is controlled by his/her mother only able to control their own monitoring field.

***

## Dataform on Database Management

### Master

```js
{
  Roles: [
    "Admin",
    "Facility",
    "Mom",
    "Child",
  ],
  MonitoringState: [
    "Not Yet",
    "Already",
  ],
  ChildMonitoringPlace: [
    "Home",
    "Facility",
  ]
}
```

### User

```js
{
  id: "id",
  type: "admin" || "facility" || "mother" || "child",
  firstname: "<First Name>",
  lastname: "<Last Name>",
  address: "<User Address>",
  identifier: "<Admin Identifier>" || "<Facility Registration Number>" || "<Mother NIK>" || "<Child NIK>",
  profilepictureurl: "<Url Profile Picture>"
}
```

### Facility Pictures

```js
{
  id: "",
  facilityid: "<User(Facility) ID>",
  url: "<Url Facility Picture>"
}
```
