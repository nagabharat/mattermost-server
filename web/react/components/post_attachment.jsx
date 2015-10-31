// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

const TextFormatting = require('../utils/text_formatting.jsx');

export default class PostAttachment extends React.Component {
    constructor(props) {
        super(props);

        this.getFieldsTable = this.getFieldsTable.bind(this);
    }

    getFieldsTable() {
        const fields = this.props.attachment.fields;
        if (!fields || !fields.length) {
            return '';
        }

        const compactTable = fields.filter((field) => field.short).length > 0;
        let tHead;
        let tBody;

        if (compactTable) {
            let headerCols = [];
            let bodyCols = [];

            fields.forEach((field, i) => {
                headerCols.push(
                    <th
                        className='attachment___field-caption'
                        key={'attachment__field-caption-' + i}
                    >
                        {field.title}
                    </th>
                );
                bodyCols.push(
                    <td
                        className='attachment___field'
                        key={'attachment__field-' + i}
                        dangerouslySetInnerHTML={{__html: TextFormatting.formatText(field.value || '')}}
                    >
                    </td>
                );
            });

            tHead = (
                <tr>
                    {headerCols}
                </tr>
            );
            tBody = (
                <tr>
                    {bodyCols}
                </tr>
            );
        } else {
            tBody = [];

            fields.forEach((field, i) => {
                tBody.push(
                    <tr key={'attachment__field-' + i}>
                        <td
                            className='attachment___field-caption'
                        >
                            {field.title}
                        </td>
                        <td
                            className='attachment___field'
                            dangerouslySetInnerHTML={{__html: TextFormatting.formatText(field.value || '')}}
                        >
                        </td>
                    </tr>
                );
            });
        }

        return (
            <table
                className='attachment___fields'
            >
                <thead>
                    {tHead}
                </thead>
                <tbody>
                    {tBody}
                </tbody>
            </table>
        );
    }

    render() {
        const data = this.props.attachment;

        let preText;
        if (data.pretext) {
            preText = (
                <div
                    className='attachment__thumb-pretext'
                    dangerouslySetInnerHTML={{__html: TextFormatting.formatText(data.pretext)}}
                >
                </div>
            );
        }

        let author = [];
        if (data.author_name || data.author_icon) {
            if (data.author_icon) {
                author.push(
                    <img
                        className='attachment__author-icon'
                        src={data.author_icon}
                        key={'attachment__author-icon'}
                        height='14'
                        width='14'
                    />
                );
            }
            if (data.author_name) {
                author.push(
                    <span
                        className='attachment__author-name'
                        key={'attachment__author-name'}
                    >
                        {data.author_name}
                    </span>
                );
            }
        }
        if (data.author_link) {
            author = (
                <a
                    href={data.author_link}
                    target='_blank'
                >
                    {author}
                </a>
            );
        }

        let title;
        if (data.title) {
            if (data.title_link) {
                title = (
                    <h1
                        className='attachment__title'
                    >
                        <a
                            className='attachment__title-link'
                            href={data.title_link}
                            target='_blank'
                        >
                            {data.title}
                        </a>
                    </h1>
                );
            } else {
                title = (
                    <h1
                        className='attachment__title'
                    >
                        {data.title}
                    </h1>
                );
            }
        }

        let text;
        if (data.text) {
            text = (
                <div
                    className='attachment__text'
                    dangerouslySetInnerHTML={{__html: TextFormatting.formatText(data.text || '')}}
                >
                </div>
            );
        }

        let image;
        if (data.image_url) {
            image = (
                <img
                    className='attachment__image'
                    src={data.image_url}
                />
            );
        }

        let thumb;
        if (data.thumb_url) {
            thumb = (
                <div
                    className='attachment__thumb-container'
                >
                    <img
                        src={data.thumb_url}
                    />
                </div>
            );
        }

        const fields = this.getFieldsTable();

        let useBorderStyle;
        if (data.color && data.color[0] === '#') {
            useBorderStyle = {borderLeftColor: data.color};
        }

        return (
            <div
                className='attachment'
            >
                {preText}
                <div className='attachment__content'>
                    <div
                        className={useBorderStyle ? 'attachment__container' : 'attachment__container attachment__container--' + data.color}
                        style={useBorderStyle}
                    >
                        {author}
                        {title}
                        <div>
                            <div
                                className={thumb ? 'attachment__body' : 'attachment__body attachment__body--no_thumb'}
                            >
                                {text}
                                {image}
                                {fields}
                            </div>
                            {thumb}
                            <div style={{clear: 'both'}}></div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}

PostAttachment.propTypes = {
    attachment: React.PropTypes.object.isRequired
};