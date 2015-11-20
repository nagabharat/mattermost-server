// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

import ChannelStore from '../stores/channel_store.jsx';
import TeamStore from '../stores/team_store.jsx';
import Constants from '../utils/constants.jsx';

function setupHomePage() {
    var last = ChannelStore.getLastVisitedName();
    if (last == null || last.length === 0) {
        window.location = TeamStore.getCurrentTeamUrl() + '/channels/' + Constants.DEFAULT_CHANNEL;
    } else {
        window.location = TeamStore.getCurrentTeamUrl() + '/channels/' + last;
    }
}

global.window.setup_home_page = setupHomePage;
