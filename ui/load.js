let id = 0;
module.exports = function(requestId) {
    return {
    	message: requestId,
		title: ""+(id++)
    };
};