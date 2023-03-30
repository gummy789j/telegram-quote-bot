package domain

type UpdateMessageResponse struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateID int64 `json:"update_id"`
		Message  *struct {
			MessageID int64 `json:"message_id"`
			From      *struct {
				ID        int64  `json:"id"`
				IsBot     bool   `json:"is_bot"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
			} `json:"from"`
			Chat *struct {
				ID                          int64  `json:"id"`
				Title                       string `json:"title"`
				Type                        string `json:"type"`
				AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
			} `json:"chat"`
			Date               int64 `json:"date"`
			NewChatParticipant *struct {
				ID        int64  `json:"id"`
				IsBot     bool   `json:"is_bot"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
			} `json:"new_chat_participant"`
			NewChatMember *struct {
				ID        int64  `json:"id"`
				IsBot     bool   `json:"is_bot"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
			} `json:"new_chat_member"`
			NewChatMembers []struct {
				ID        int64  `json:"id"`
				IsBot     bool   `json:"is_bot"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
			} `json:"new_chat_members"`
			Text     *string `json:"text"`
			Entities []struct {
				Offset int64  `json:"offset"`
				Length int64  `json:"length"`
				Type   string `json:"type"`
			} `json:"entities"`
		} `json:"message"`
	} `json:"result"`
}

/*
{
    "ok": true,
    "result": [
        {
            "update_id": 926617501,
            "message": {
                "message_id": 700,
                "from": {
                    "id": 6040823283,
                    "is_bot": false,
                    "first_name": "Hsu",
                    "last_name": "Sunny"
                },
                "chat": {
                    "id": -781207517,
                    "title": "MaxThingsRyght",
                    "type": "group",
                    "all_members_are_administrators": true
                },
                "date": 1680149700,
                "new_chat_participant": {
                    "id": 6040823283,
                    "is_bot": false,
                    "first_name": "Hsu",
                    "last_name": "Sunny"
                },
                "new_chat_member": {
                    "id": 6040823283,
                    "is_bot": false,
                    "first_name": "Hsu",
                    "last_name": "Sunny"
                },
                "new_chat_members": [
                    {
                        "id": 6040823283,
                        "is_bot": false,
                        "first_name": "Hsu",
                        "last_name": "Sunny"
                    }
                ]
            }
        },
        {
            "update_id": 926617502,
            "message": {
                "message_id": 727,
                "from": {
                    "id": 6094106269,
                    "is_bot": false,
                    "first_name": "Jason",
                    "last_name": "Teng"
                },
                "chat": {
                    "id": -781207517,
                    "title": "MaxThingsRyght",
                    "type": "group",
                    "all_members_are_administrators": true
                },
                "date": 1680151498,
                "new_chat_participant": {
                    "id": 1082588647,
                    "is_bot": false,
                    "first_name": "Ann",
                    "username": "anneammy"
                },
                "new_chat_member": {
                    "id": 1082588647,
                    "is_bot": false,
                    "first_name": "Ann",
                    "username": "anneammy"
                },
                "new_chat_members": [
                    {
                        "id": 1082588647,
                        "is_bot": false,
                        "first_name": "Ann",
                        "username": "anneammy"
                    }
                ]
            }
        },
        {
            "update_id": 926617503,
            "message": {
                "message_id": 946,
                "from": {
                    "id": 1881712391,
                    "is_bot": false,
                    "first_name": "Steven",
                    "last_name": "Lin",
                    "username": "gummy789j",
                    "language_code": "zh-hans"
                },
                "chat": {
                    "id": 1881712391,
                    "first_name": "Steven",
                    "last_name": "Lin",
                    "username": "gummy789j",
                    "type": "private"
                },
                "date": 1680170995,
                "text": "/notify",
                "entities": [
                    {
                        "offset": 0,
                        "length": 7,
                        "type": "bot_command"
                    }
                ]
            }
        },
        {
            "update_id": 926617504,
            "message": {
                "message_id": 947,
                "from": {
                    "id": 1881712391,
                    "is_bot": false,
                    "first_name": "Steven",
                    "last_name": "Lin",
                    "username": "gummy789j",
                    "language_code": "zh-hans"
                },
                "chat": {
                    "id": 1881712391,
                    "first_name": "Steven",
                    "last_name": "Lin",
                    "username": "gummy789j",
                    "type": "private"
                },
                "date": 1680171868,
                "text": "/hey",
                "entities": [
                    {
                        "offset": 0,
                        "length": 4,
                        "type": "bot_command"
                    }
                ]
            }
        },
        {
            "update_id": 926617505,
            "message": {
                "message_id": 948,
                "from": {
                    "id": 1881712391,
                    "is_bot": false,
                    "first_name": "Steven",
                    "last_name": "Lin",
                    "username": "gummy789j",
                    "language_code": "en"
                },
                "chat": {
                    "id": -781207517,
                    "title": "MaxThingsRyght",
                    "type": "group",
                    "all_members_are_administrators": false
                },
                "date": 1680173682,
                "new_chat_participant": {
                    "id": 5050867637,
                    "is_bot": false,
                    "first_name": "Jenny",
                    "username": "jennystagramm"
                },
                "new_chat_member": {
                    "id": 5050867637,
                    "is_bot": false,
                    "first_name": "Jenny",
                    "username": "jennystagramm"
                },
                "new_chat_members": [
                    {
                        "id": 5050867637,
                        "is_bot": false,
                        "first_name": "Jenny",
                        "username": "jennystagramm"
                    }
                ]
            }
        },
        {
            "update_id": 926617506,
            "my_chat_member": {
                "chat": {
                    "id": -905284654,
                    "title": "test-bot",
                    "type": "group",
                    "all_members_are_administrators": true
                },
                "from": {
                    "id": 1881712391,
                    "is_bot": false,
                    "first_name": "Steven",
                    "last_name": "Lin",
                    "username": "gummy789j",
                    "language_code": "zh-hans"
                },
                "date": 1680182231,
                "old_chat_member": {
                    "user": {
                        "id": 6156662592,
                        "is_bot": true,
                        "first_name": "GreenHat",
                        "username": "gummy_s_bot"
                    },
                    "status": "left"
                },
                "new_chat_member": {
                    "user": {
                        "id": 6156662592,
                        "is_bot": true,
                        "first_name": "GreenHat",
                        "username": "gummy_s_bot"
                    },
                    "status": "member"
                }
            }
        },
        {
            "update_id": 926617507,
            "message": {
                "message_id": 949,
                "from": {
                    "id": 1881712391,
                    "is_bot": false,
                    "first_name": "Steven",
                    "last_name": "Lin",
                    "username": "gummy789j",
                    "language_code": "zh-hans"
                },
                "chat": {
                    "id": -905284654,
                    "title": "test-bot",
                    "type": "group",
                    "all_members_are_administrators": true
                },
                "date": 1680182231,
                "new_chat_participant": {
                    "id": 6156662592,
                    "is_bot": true,
                    "first_name": "GreenHat",
                    "username": "gummy_s_bot"
                },
                "new_chat_member": {
                    "id": 6156662592,
                    "is_bot": true,
                    "first_name": "GreenHat",
                    "username": "gummy_s_bot"
                },
                "new_chat_members": [
                    {
                        "id": 6156662592,
                        "is_bot": true,
                        "first_name": "GreenHat",
                        "username": "gummy_s_bot"
                    }
                ]
            }
        },
        {
            "update_id": 926617508,
            "message": {
                "message_id": 950,
                "from": {
                    "id": 1881712391,
                    "is_bot": false,
                    "first_name": "Steven",
                    "last_name": "Lin",
                    "username": "gummy789j",
                    "language_code": "zh-hans"
                },
                "chat": {
                    "id": -905284654,
                    "title": "test-bot",
                    "type": "group",
                    "all_members_are_administrators": true
                },
                "date": 1680182238,
                "text": "/help",
                "entities": [
                    {
                        "offset": 0,
                        "length": 5,
                        "type": "bot_command"
                    }
                ]
            }
        }
    ]
}
*/
