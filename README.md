# Azure Logic Apps Demo

# !!! WARNING !!!

There is zero security in this demo. All of the resources are deployed into their own little world for this reason but you should be aware that 
we might as well have passwords hard-coded :)

# Purpose

This demo is meant to mimic a [Logic App](https://azure.microsoft.com/en-us/services/logic-apps/) that is used as a way to funnel a single request to various 3rd party/external APIs. The concept is that you have a single HTTP endpoint (the Logic App) that you can make requests against. The Logic App will take that request and send it to other 3rd party APIs. The results from those 3rd party APIs are then used inside of the Logic App to update a single "master of record" type database.

In this case we are mocking up a conference speaker submission app. The speaker's proposed talk will be submitted in a `PENDING` state. This submission will then be handed off to mock 3rd party APIs that "validate" the submission and update the state to `CONFIRMED`.

The mock 3rd party APIs are handled by a simple [Go](https://golang.org/) program that is located inside of this repo.

# Deployment

1. Clone this repo :)
2. Deploy ARM template
    1. Login to the Azure portal
    2. Click "Create a resource" in the top left corner
    3. In the search box type `Template Deployment` and select it
    4. Click "Create"
    5. Click "Build your own template in the editor"
    6. Click "Load File" and select the `template.json` file from the `templates` directory in this repo
    7. Click "Save"
    8. Fill out the parameters with the values you want. Take note of the parameters that require unique values, mainly the DNS prefix
    9. Sit back while things deploy
3. Once the deployment is complete you can go to the resource group that you specified in step #2 and you should see all of the various resources used in this demo
   1. Click on the Logic App
   2. Click on `Edit` or `Logic App Designer`
   3. The top of the flow will be the endpoint for testing. Click on it and copy the URL. This is auto-generated and will change for each deployment
4. Take the URL from step #3 and load it into your favorite REST testing tool ([Postman](https://www.getpostman.com/) for example)
    1. Set the `Content-Type` header to `application/json`
    1. Set the body to:
   
        ```json
        {
	        "name":"Demo Name",
	        "email": "demo@demo.com",
	        "topic": "Demo Topic",
	        "status": "PENDING"
        }
        ```
    1. POST it!

5. Now get your VMs IP and verify
    1. Back in the resource group click on the VM
    2. Take note of the `Public IP Address`
    3. Visit `http://IP_from_above_step:8000`
    4. You should see a basic page that shows various entries. Notice that the one you just submitted has a `status` of `COMPLETE`. This was handled in the Logic Apps workflow by simulating 3rd party APIs
