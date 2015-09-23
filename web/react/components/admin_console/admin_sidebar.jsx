// Copyright (c) 2015 Spinpunch, Inc. All Rights Reserved.
// See License.txt for license information.

var AdminSidebarHeader = require('./admin_sidebar_header.jsx');

export default class AdminSidebar extends React.Component {
    constructor(props) {
        super(props);

        this.isSelected = this.isSelected.bind(this);
        this.handleClick = this.handleClick.bind(this);

        this.state = {
        };
    }

    handleClick(name, e) {
        e.preventDefault();
        this.props.selectTab(name);
    }

    isSelected(name) {
        if (this.props.selected === name) {
            return 'active';
        }

        return '';
    }

    componentDidMount() {
    }

    render() {
        return (
            <div className='sidebar--left sidebar--collapsable'>
                <div>
                    <AdminSidebarHeader />
                    <ul className='nav nav-pills nav-stacked'>
                        <li>
                            <ul className='nav nav__sub-menu'>
                                <li>
                                    <a
                                        href='#'
                                        className={this.isSelected('service_settings')}
                                        onClick={this.handleClick.bind(this, 'service_settings')}
                                    >
                                        {'Service Settings'}
                                    </a>
                                </li>
                                <li>
                                    <a
                                        href='#'
                                        className={this.isSelected('team_settings')}
                                        onClick={this.handleClick.bind(this, 'team_settings')}
                                    >
                                        {'Team Settings'}
                                    </a>
                                </li>
                                <li>
                                    <a
                                        href='#'
                                        className={this.isSelected('sql_settings')}
                                        onClick={this.handleClick.bind(this, 'sql_settings')}
                                    >
                                        {'SQL Settings'}
                                    </a>
                                </li>
                                <li>
                                    <a
                                        href='#'
                                        className={this.isSelected('email_settings')}
                                        onClick={this.handleClick.bind(this, 'email_settings')}
                                    >
                                        {'Email Settings'}
                                    </a>
                                </li>
                                <li>
                                    <a
                                        href='#'
                                        className={this.isSelected('image_settings')}
                                        onClick={this.handleClick.bind(this, 'image_settings')}
                                    >
                                        {'File Settings'}
                                    </a>
                                </li>
                                <li>
                                    <a
                                        href='#'
                                        className={this.isSelected('log_settings')}
                                        onClick={this.handleClick.bind(this, 'log_settings')}
                                    >
                                        {'Log Settings'}
                                    </a>
                                </li>
                                <li>
                                    <a
                                        href='#'
                                        className={this.isSelected('logs')}
                                        onClick={this.handleClick.bind(this, 'logs')}
                                    >
                                        {'Logs'}
                                    </a>
                                </li>
                                <li>
                                    <a
                                        href='#'
                                        className={this.isSelected('rate_settings')}
                                        onClick={this.handleClick.bind(this, 'rate_settings')}
                                    >
                                        {'Rate Limit Settings'}
                                    </a>
                                </li>
                                <li>
                                    <a
                                        href='#'
                                        className={this.isSelected('privacy_settings')}
                                        onClick={this.handleClick.bind(this, 'privacy_settings')}
                                    >
                                        {'Privacy Settings'}
                                    </a>
                                </li>
                                <li>
                                    <a
                                        href='#'
                                        className={this.isSelected('gitlab_settings')}
                                        onClick={this.handleClick.bind(this, 'gitlab_settings')}
                                    >
                                        {'GitLab Settings'}
                                    </a>
                                </li>
                            </ul>
                        </li>
                    </ul>
                </div>
            </div>
        );
    }
}

AdminSidebar.propTypes = {
    selected: React.PropTypes.string,
    selectTab: React.PropTypes.func
};