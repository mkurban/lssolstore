# lssproto

Lotsize Solution Store Protocol

	-> Read Problem in: (Problem[JSON]) => (ProbID, Err)
	-> Create Solution (ProbID) => (SolID, Err)

	-> Get Problem Info (ProbID) => (, Err)
	-> Get Bucket Info (ProbID, BucketID) => (, Err)

	-> Start Improving Solution (SolID) => (ImpToken, Err)
		+=> Improvement started for Solution {SolID}
	-> Add Cols To Solution  (SolID, ImpToken, [Cols]) => Err
		+=> Cols added to Solution {SolID} ([Cols])
	-> Get Cols In Solution  (SolID) => ([Cols], Err)
	-> Stop Improving Solution (SolID, ImpToken, [Selected Cols]) => Err
		+=> Improvement stopped for Solution {SolID} (Reason)
