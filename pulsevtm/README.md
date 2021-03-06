# Terraform Pulse Virtual Traffic Manager Provider

This document is in place for developer documentation.  User documentation is located [HERE](https://www.terraform.io/docs/providers/pulsevtm/) on Terraform's website.

A Terraform provider for the Pulse vTM.  The pulse vTM provider is used to interact with resources supported by the pulse Virtual Traffic Manager (vTM).
The provider needs to be configured with the proper credentials before it can be used.

## Introductory Documentation

Both [README.md](../../../README.md) and [BUILDING.md](../../../BUILDING.md) should be read first!

## Base API Dependency ~ [go-pulse-vtm](https://github.com/sky-uk/go-pulse-vtm)

This provider utilizes [go-pulse-vtm](https://github.com/sky-uk/go-pulse-vtm) Go Library for communicating to the Pulse Virtual Traffic Manager REST API.
Because of the dependency this provider is compatible with Pulse systems that are supported by go-pulse-vtm. If you want to contributed additional functionality into gopulse-vtm API bindings
please feel free to send the pull requests.


## Resources Implemented
| Feature                 | Create | Read  | Update  | Delete |
|-------------------------|--------|-------|---------|--------|
| Monitor                 |   Y    |   Y   |    N    |   Y    |
| Pools                   |   N    |   N   |    N    |   N    |
| Traffic IP              |   Y    |   Y   |    Y    |   Y    |
| Virtual Server          |   N    |   N   |    N    |   N    |


### Traffic IP Group

Implemented Traffic IP Group Attributes  
  
| Attribute       | Create | Read | Update | Delete*** |  
|-----------------|--------|------|--------|-----------|  
| name            |    Y   |   Y  |   N*   |   Y       |  
| enabled         |    Y   |   Y  |   Y    |   Y       |  
| hashsourceport  |    Y   |   Y  |   Y    |   Y       |  
| ipaddresses     |    Y   |   Y  |   Y    |   Y       |  
| trafficmanagers |    Y** |   Y  |   Y**  |   Y       |  
| mode            |    Y   |   Y  |   Y    |   Y       |  
| multicastip     |    Y   |   Y  |   Y    |   Y       |  
  
   
*Changing the name attribute will force (delete/create) the creation of a new resource.  
**trafficmanagers is dynamically generated on create and update. Its run each time an update is applied.  
***Attributes are deleted when the resource is removed.  

### Sample Traffic IP Group template

```   
resource "pulsevtm_traffic_ip_group" "my_app" {
  name = "sample-traffic-ip-group"
  enabled = true
  hashsourceport = true
  ipaddresses = ["10.34.10.12]
  multicastip = "239.191.130.3"
  mode = "multihosted"
}
```   

*Note: The trafficmanagers attribute is automatically generated by the create and update functions and can't be user specified.  
   
### Limitations

This is currently a proof of concept and only has a very limited number of
supported resources.  These resources also have a very limited number
of attributes.

This section is a work in progress and additional contributions are more than welcome.
