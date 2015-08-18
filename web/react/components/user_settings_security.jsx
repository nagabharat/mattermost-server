var SettingItemMin = require('./setting_item_min.jsx');
var SettingItemMax = require('./setting_item_max.jsx');
var client = require('../utils/client.jsx');
var AsyncClient = require('../utils/async_client.jsx');
var Constants = require('../utils/constants.jsx');

module.exports = React.createClass({
    displayName: 'SecurityTab',
    submitPassword: function(e) {
        e.preventDefault();

        var user = this.props.user;
        var currentPassword = this.state.currentPassword;
        var newPassword = this.state.newPassword;
        var confirmPassword = this.state.confirmPassword;

        if (currentPassword === '') {
            this.setState({passwordError: 'Please enter your current password', serverError: ''});
            return;
        }

        if (newPassword.length < 5) {
            this.setState({passwordError: 'New passwords must be at least 5 characters', serverError: ''});
            return;
        }

        if (newPassword !== confirmPassword) {
            this.setState({passwordError: 'The new passwords you entered do not match', serverError: ''});
            return;
        }

        var data = {};
        data.user_id = user.id;
        data.current_password = currentPassword;
        data.new_password = newPassword;

        client.updatePassword(data,
            function() {
                this.props.updateSection('');
                AsyncClient.getMe();
                this.setState({currentPassword: '', newPassword: '', confirmPassword: ''});
            }.bind(this),
            function(err) {
                var state = this.getInitialState();
                if (err.message) {
                    state.serverError = err.message;
                } else {
                    state.serverError = err;
                }
                state.passwordError = '';
                this.setState(state);
            }.bind(this)
        );
    },
    updateCurrentPassword: function(e) {
        this.setState({currentPassword: e.target.value});
    },
    updateNewPassword: function(e) {
        this.setState({newPassword: e.target.value});
    },
    updateConfirmPassword: function(e) {
        this.setState({confirmPassword: e.target.value});
    },
    handleHistoryOpen: function() {
        this.setState({willReturn: true});
        $("#user_settings").modal('hide');
    },
    handleDevicesOpen: function() {
        this.setState({willReturn: true});
        $("#user_settings").modal('hide');
    },
    handleClose: function() {
        $(this.getDOMNode()).find('.form-control').each(function() {
            this.value = '';
        });
        this.setState({currentPassword: '', newPassword: '', confirmPassword: '', serverError: null, passwordError: null});

        if (!this.state.willReturn) {
            this.props.updateTab('general');
        } else {
            this.setState({willReturn: false});
        }
    },
    componentDidMount: function() {
        $('#user_settings').on('hidden.bs.modal', this.handleClose);
    },
    componentWillUnmount: function() {
        $('#user_settings').off('hidden.bs.modal', this.handleClose);
        this.props.updateSection('');
    },
    getInitialState: function() {
        return {currentPassword: '', newPassword: '', confirmPassword: '', willReturn: false};
    },
    render: function() {
        var serverError = this.state.serverError ? this.state.serverError : null;
        var passwordError = this.state.passwordError ? this.state.passwordError : null;

        var updateSectionStatus;
        var passwordSection;
        var self = this;
        if (this.props.activeSection === 'password') {
            var inputs = [];
            var submit = null;

            if (this.props.user.auth_service === '') {
                inputs.push(
                    <div className='form-group'>
                        <label className='col-sm-5 control-label'>Current Password</label>
                        <div className='col-sm-7'>
                            <input className='form-control' type='password' onChange={this.updateCurrentPassword} value={this.state.currentPassword}/>
                        </div>
                    </div>
                );
                inputs.push(
                    <div className='form-group'>
                        <label className='col-sm-5 control-label'>New Password</label>
                        <div className='col-sm-7'>
                            <input className='form-control' type='password' onChange={this.updateNewPassword} value={this.state.newPassword}/>
                        </div>
                    </div>
                );
                inputs.push(
                    <div className='form-group'>
                        <label className='col-sm-5 control-label'>Retype New Password</label>
                        <div className='col-sm-7'>
                            <input className='form-control' type='password' onChange={this.updateConfirmPassword} value={this.state.confirmPassword}/>
                        </div>
                    </div>
                );

                submit = this.submitPassword;
            } else {
                inputs.push(
                    <div className='form-group'>
                        <label className='col-sm-12'>Log in occurs through GitLab. Please see your GitLab account settings page to update your password.</label>
                    </div>
                );
            }

            updateSectionStatus = function(e) {
                self.props.updateSection('');
                self.setState({currentPassword: '', newPassword: '', confirmPassword: '', serverError: null, passwordError: null});
                e.preventDefault();
            };

            passwordSection = (
                <SettingItemMax
                    title='Password'
                    inputs={inputs}
                    submit={submit}
                    server_error={serverError}
                    client_error={passwordError}
                    updateSection={updateSectionStatus}
                />
            );
        } else {
            var describe;
            if (this.props.user.auth_service === '') {
                var d = new Date(this.props.user.last_password_update);
                var hour = d.getHours() % 12 ? String(d.getHours() % 12) : '12';
                var min = d.getMinutes() < 10 ? '0' + d.getMinutes() : String(d.getMinutes());
                var timeOfDay = d.getHours() >= 12 ? ' pm' : ' am';
                describe = 'Last updated ' + Constants.MONTHS[d.getMonth()] + ' ' + d.getDate() + ', ' + d.getFullYear() + ' at ' + hour + ':' + min + timeOfDay;
            } else {
                describe = 'Log in done through GitLab';
            }

            updateSectionStatus = function() {
                self.props.updateSection('password');
            };

            passwordSection = (
                <SettingItemMin
                    title='Password'
                    describe={describe}
                    updateSection={updateSectionStatus}
                />
            );
        }

        return (
            <div>
                <div className='modal-header'>
                    <button type='button' className='close' data-dismiss='modal' aria-label='Close'><span aria-hidden='true'>&times;</span></button>
                    <h4 className='modal-title' ref='title'><i className='modal-back'></i>Security Settings</h4>
                </div>
                <div className='user-settings'>
                    <h3 className='tab-header'>Security Settings</h3>
                    <div className='divider-dark first'/>
                    {passwordSection}
                    <div className='divider-dark'/>
                    <br></br>
                    <a data-toggle='modal' className='security-links theme' data-target='#access-history' href='#' onClick={this.handleHistoryOpen}><i className='fa fa-clock-o'></i>View Access History</a>
                    <b> </b>
                    <a data-toggle='modal' className='security-links theme' data-target='#activity-log' href='#' onClick={this.handleDevicesOpen}><i className='fa fa-globe'></i>View and Logout of Active Sessions</a>
                </div>
            </div>
        );
    }
});
