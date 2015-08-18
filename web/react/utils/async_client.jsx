// Copyright (c) 2015 Spinpunch, Inc. All Rights Reserved.
// See License.txt for license information.

var client = require('./client.jsx');
var AppDispatcher = require('../dispatcher/app_dispatcher.jsx');
var ChannelStore = require('../stores/channel_store.jsx');
var ConfigStore = require('../stores/config_store.jsx');
var PostStore = require('../stores/post_store.jsx');
var UserStore = require('../stores/user_store.jsx');
var utils = require('./utils.jsx');

var Constants = require('./constants.jsx');
var ActionTypes = Constants.ActionTypes;

// Used to track in progress async calls
var callTracker = {};

function dispatchError(err, method) {
    AppDispatcher.handleServerAction({
        type: ActionTypes.RECIEVED_ERROR,
        err: err,
        method: method
    });
}
module.exports.dispatchError = dispatchError;

function isCallInProgress(callName) {
    if (!(callName in callTracker)) {
        return false;
    }

    if (callTracker[callName] === 0) {
        return false;
    }

    if (utils.getTimestamp() - callTracker[callName] > 5000) {
        console.log('AsyncClient call ' + callName + ' expired after more than 5 seconds');
        return false;
    }

    return true;
}

function getChannels(force, updateLastViewed, checkVersion) {
    var channels = ChannelStore.getAll();

    if (channels.length === 0 || force) {
        if (isCallInProgress('getChannels')) {
            return;
        }

        callTracker.getChannels = utils.getTimestamp();

        client.getChannels(
            function(data, textStatus, xhr) {
                callTracker.getChannels = 0;

                if (checkVersion) {
                    var serverVersion = xhr.getResponseHeader('X-Version-ID');

                    if (!UserStore.getLastVersion()) {
                        UserStore.setLastVersion(serverVersion);
                    }

                    if (serverVersion !== UserStore.getLastVersion()) {
                        UserStore.setLastVersion(serverVersion);
                        window.location.href = window.location.href;
                        console.log('Detected version update refreshing the page');
                    }
                }

                if (xhr.status === 304 || !data) {
                    return;
                }

                AppDispatcher.handleServerAction({
                    type: ActionTypes.RECIEVED_CHANNELS,
                    channels: data.channels,
                    members: data.members
                });
            },
            function(err) {
                callTracker.getChannels = 0;
                dispatchError(err, 'getChannels');
            }
        );
    } else {
        if (isCallInProgress('getChannelCounts')) {
            return;
        }

        callTracker.getChannelCounts = utils.getTimestamp();

        client.getChannelCounts(
            function(data, textStatus, xhr) {
                callTracker.getChannelCounts = 0;

                if (xhr.status === 304 || !data) {
                    return;
                }

                var countMap = data.counts;
                var updateAtMap = data.update_times;

                for (var id in countMap) {
                    var c = ChannelStore.get(id);
                    var count = countMap[id];
                    var updateAt = updateAtMap[id];
                    if (!c || c.total_msg_count !== count || updateAt > c.update_at) {
                        getChannel(id);
                    }
                }
            },
            function(err) {
                callTracker.getChannelCounts = 0;
                dispatchError(err, 'getChannelCounts');
            }
        );
    }

    if (updateLastViewed && ChannelStore.getCurrentId() != null) {
        module.exports.updateLastViewedAt();
    }
}
module.exports.getChannels = getChannels;

function getChannel(id) {
    if (isCallInProgress('getChannel' + id)) {
        return;
    }

    callTracker['getChannel' + id] = utils.getTimestamp();

    client.getChannel(id,
        function(data, textStatus, xhr) {
            callTracker['getChannel' + id] = 0;

            if (xhr.status === 304 || !data) {
                return;
            }

            AppDispatcher.handleServerAction({
                type: ActionTypes.RECIEVED_CHANNEL,
                channel: data.channel,
                member: data.member
            });
        },
        function(err) {
            callTracker['getChannel' + id] = 0;
            dispatchError(err, 'getChannel');
        }
    );
}
module.exports.getChannel = getChannel;

module.exports.updateLastViewedAt = function() {
    if (isCallInProgress('updateLastViewed')) return;

    if (ChannelStore.getCurrentId() == null) return;

    callTracker['updateLastViewed'] = utils.getTimestamp();
    client.updateLastViewedAt(
        ChannelStore.getCurrentId(),
        function(data) {
            callTracker['updateLastViewed'] = 0;
        },
        function(err) {
            callTracker['updateLastViewed'] = 0;
            dispatchError(err, 'updateLastViewedAt');
        }
    );
}

module.exports.getMoreChannels = function(force) {
    if (isCallInProgress('getMoreChannels')) return;

    if (ChannelStore.getMoreAll().loading || force) {

        callTracker['getMoreChannels'] = utils.getTimestamp();
        client.getMoreChannels(
            function(data, textStatus, xhr) {
                callTracker['getMoreChannels'] = 0;

                if (xhr.status === 304 || !data) return;

                AppDispatcher.handleServerAction({
                    type: ActionTypes.RECIEVED_MORE_CHANNELS,
                    channels: data.channels,
                    members: data.members
                });
            },
            function(err) {
                callTracker['getMoreChannels'] = 0;
                dispatchError(err, 'getMoreChannels');
            }
        );
    }
}

module.exports.getChannelExtraInfo = function(force) {
    var channelId = ChannelStore.getCurrentId();

    if (channelId != null) {
        if (isCallInProgress('getChannelExtraInfo_'+channelId)) return;
        var minMembers = ChannelStore.getCurrent() && ChannelStore.getCurrent().type === 'D' ? 1 : 0;

        if (ChannelStore.getCurrentExtraInfo().members.length <= minMembers || force) {
            callTracker['getChannelExtraInfo_'+channelId] = utils.getTimestamp();
            client.getChannelExtraInfo(
                channelId,
                function(data, textStatus, xhr) {
                    callTracker['getChannelExtraInfo_'+channelId] = 0;

                    if (xhr.status === 304 || !data) return;

                    AppDispatcher.handleServerAction({
                        type: ActionTypes.RECIEVED_CHANNEL_EXTRA_INFO,
                        extra_info: data
                    });
                },
                function(err) {
                    callTracker['getChannelExtraInfo_'+channelId] = 0;
                    dispatchError(err, 'getChannelExtraInfo');
                }
            );
        }
    }
}

module.exports.getProfiles = function() {
    if (isCallInProgress('getProfiles')) return;

    callTracker['getProfiles'] = utils.getTimestamp();
    client.getProfiles(
        function(data, textStatus, xhr) {
            callTracker['getProfiles'] = 0;

            if (xhr.status === 304 || !data) return;

            AppDispatcher.handleServerAction({
                type: ActionTypes.RECIEVED_PROFILES,
                profiles: data
            });
        },
        function(err) {
            callTracker['getProfiles'] = 0;
            dispatchError(err, 'getProfiles');
        }
    );
}

module.exports.getSessions = function() {
    if (isCallInProgress('getSessions')) return;

    callTracker['getSessions'] = utils.getTimestamp();
    client.getSessions(
        UserStore.getCurrentId(),
        function(data, textStatus, xhr) {
            callTracker['getSessions'] = 0;

            if (xhr.status === 304 || !data) return;

            AppDispatcher.handleServerAction({
                type: ActionTypes.RECIEVED_SESSIONS,
                sessions: data
            });
        },
        function(err) {
            callTracker['getSessions'] = 0;
            dispatchError(err, 'getSessions');
        }
    );
}

module.exports.getAudits = function() {
    if (isCallInProgress('getAudits')) return;

    callTracker['getAudits'] = utils.getTimestamp();
    client.getAudits(
        UserStore.getCurrentId(),
        function(data, textStatus, xhr) {
            callTracker['getAudits'] = 0;

            if (xhr.status === 304 || !data) return;

            AppDispatcher.handleServerAction({
                type: ActionTypes.RECIEVED_AUDITS,
                audits: data
            });
        },
        function(err) {
            callTracker['getAudits'] = 0;
            dispatchError(err, 'getAudits');
        }
    );
}

module.exports.findTeams = function(email) {
    if (isCallInProgress('findTeams_'+email)) return;

    var user = UserStore.getCurrentUser();
    if (user) {
        callTracker['findTeams_'+email] = utils.getTimestamp();
        client.findTeams(
            user.email,
            function(data, textStatus, xhr) {
                callTracker['findTeams_'+email] = 0;

                if (xhr.status === 304 || !data) return;

                AppDispatcher.handleServerAction({
                    type: ActionTypes.RECIEVED_TEAMS,
                    teams: data
                });
            },
            function(err) {
                callTracker['findTeams_'+email] = 0;
                dispatchError(err, 'findTeams');
            }
        );
    }
}

module.exports.search = function(terms) {
    if (isCallInProgress('search_'+String(terms))) return;

    callTracker['search_'+String(terms)] = utils.getTimestamp();
    client.search(
        terms,
        function(data, textStatus, xhr) {
            callTracker['search_'+String(terms)] = 0;

            if (xhr.status === 304 || !data) return;

            AppDispatcher.handleServerAction({
                type: ActionTypes.RECIEVED_SEARCH,
                results: data
            });
        },
        function(err) {
            callTracker['search_'+String(terms)] = 0;
            dispatchError(err, 'search');
        }
    );
}

module.exports.getPosts = function(force, id, maxPosts) {
    if (PostStore.getCurrentPosts() == null || force) {
        var channelId = id;
        if (channelId == null) {
            channelId = ChannelStore.getCurrentId();
        }

        if (isCallInProgress('getPosts_' + channelId)) {
            return;
        }

        var postList = PostStore.getCurrentPosts();

        var max = maxPosts;
        if (max == null) {
            max = Constants.POST_CHUNK_SIZE * Constants.MAX_POST_CHUNKS;
        }

        // if we already have more than POST_CHUNK_SIZE posts,
        //   let's get the amount we have but rounded up to next multiple of POST_CHUNK_SIZE,
        //   with a max at maxPosts
        var numPosts = Math.min(max, Constants.POST_CHUNK_SIZE);
        if (postList && postList.order.length > 0) {
            numPosts = Math.min(max, Constants.POST_CHUNK_SIZE * Math.ceil(postList.order.length / Constants.POST_CHUNK_SIZE));
        }

        if (channelId != null) {
            callTracker['getPosts_' + channelId] = utils.getTimestamp();

            client.getPosts(
                channelId,
                0,
                numPosts,
                function(data, textStatus, xhr) {
                    if (xhr.status === 304 || !data) return;

                    AppDispatcher.handleServerAction({
                        type: ActionTypes.RECIEVED_POSTS,
                        id: channelId,
                        post_list: data
                    });

                    module.exports.getProfiles();
                },
                function(err) {
                    dispatchError(err, 'getPosts');
                },
                function() {
                    callTracker['getPosts_' + channelId] = 0;
                }
            );
        }
    }
}

function getMe() {
    if (isCallInProgress('getMe')) {
        return;
    }

    callTracker.getMe = utils.getTimestamp();
    client.getMeSynchronous(
        function(data, textStatus, xhr) {
            callTracker.getMe = 0;

            if (xhr.status === 304 || !data) return;

            AppDispatcher.handleServerAction({
                type: ActionTypes.RECIEVED_ME,
                me: data
            });
        },
        function(err) {
            callTracker.getMe = 0;
            dispatchError(err, 'getMe');
        }
    );
}
module.exports.getMe = getMe;

module.exports.getStatuses = function() {
    if (isCallInProgress('getStatuses')) return;

    callTracker['getStatuses'] = utils.getTimestamp();
    client.getStatuses(
        function(data, textStatus, xhr) {
            callTracker['getStatuses'] = 0;

            if (xhr.status === 304 || !data) return;

            AppDispatcher.handleServerAction({
                type: ActionTypes.RECIEVED_STATUSES,
                statuses: data
            });
        },
        function(err) {
            callTracker['getStatuses'] = 0;
            dispatchError(err, 'getStatuses');
        }
    );
}

module.exports.getMyTeam = function() {
    if (isCallInProgress('getMyTeam')) return;

    callTracker['getMyTeam'] = utils.getTimestamp();
    client.getMyTeam(
        function(data, textStatus, xhr) {
            callTracker['getMyTeam'] = 0;

            if (xhr.status === 304 || !data) return;

            AppDispatcher.handleServerAction({
                type: ActionTypes.RECIEVED_TEAM,
                team: data
            });
        },
        function(err) {
            callTracker['getMyTeam'] = 0;
            dispatchError(err, 'getMyTeam');
        }
    );
}

function getConfig() {
    if (isCallInProgress('getConfig')) {
        return;
    }

    callTracker['getConfig'] = utils.getTimestamp();
    client.getConfig(
        function(data, textStatus, xhr) {
            callTracker['getConfig'] = 0;

            if (data && xhr.status !== 304) {
                AppDispatcher.handleServerAction({
                    type: ActionTypes.RECIEVED_CONFIG,
                    settings: data
                });
            }
        },
        function(err) {
            callTracker['getConfig'] = 0;
            dispatchError(err, 'getConfig');
        }
    );
}
module.exports.getConfig = getConfig;
