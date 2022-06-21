# Momento::SimpleCache::Cache

Creates a Momento serverless cache.

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "Type" : "Momento::SimpleCache::Cache",
    "Properties" : {
        "<a href="#name" title="Name">Name</a>" : <i>String</i>
    }
}
</pre>

### YAML

<pre>
Type: Momento::SimpleCache::Cache
Properties:
    <a href="#name" title="Name">Name</a>: <i>String</i>
</pre>

## Properties

#### Name

Name of the cache to be created.

_Required_: Yes

_Type_: String

_Minimum_: <code>3</code>

_Maximum_: <code>255</code>

_Pattern_: <code>^[a-zA-Z0-9-_.]{3,255}$</code>

_Update requires_: [Replacement](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-replacement)

## Return Values

### Ref

When you pass the logical ID of this resource to the intrinsic `Ref` function, Ref returns the Name.
