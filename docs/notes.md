# Design Notes

### Metrics

#### Why are there 2 metrics for positive and failed event?
The decision to have 2 different metrics for successful and failed operations was taken to reduce the risk of having metrics
with high cardinality.
> If this makes for a worst user experience overall, this approach can be changed.

#### Why are all the metrics name ending in `total`?
That is because most metrics exported are counter metrics and most counter metrics usually end in `_total` or `_count` by convention.
> If this makes for a worst user experience overall, this approach can be changed.
