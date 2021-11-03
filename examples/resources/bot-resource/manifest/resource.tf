variable bot_description {
  type = string
  description = "The description of the bot"
}

locals {

  intents = [
    {   
        id = "gas-leak"
        questions = ["I smell gas in my house. What should I do?",
                    "help my gas is leakings",
                    "emergency gas leak"]

        answer = "For Gas Emergencies or Safety Issues call Emergencies: 911 For general safety issues: 1-800-427-2200"
    },
    {
        id = "password-reset"
        questions = ["how do I reset my password",
                    "I forgot my password",
                    "Can't remember my password",
                    "My login does not work"]

        answer = "If you forgot your My Account password, securely reset it with an authorization code that is sent to your cellphone number on your My Account profile."
    }
  ]

  slot_type_values = flatten([for s in local.intents: [for q in s.questions: {"sampleValue": {"value": q}, "synonyms": []}]])  
}

data "template_file" "slot_types" {
  template = file("${path.module}/QnABot/BotLocales/en_US/SlotTypes/QnaSlotType/SlotType.json.tmpl")
  vars = {
    slot_types = jsonencode(local.slot_type_values)
  }
}

data "template_file" "bot_json" {
  template = file("${path.module}/QnABot/Bot.json.tmpl")
  vars = {
    bot_description = var.bot_description
  }
}

data "archive_file" "bot" {
  type        = "zip"
  output_path = "${path.module}/archive/bot.zip"

  # content that is templated
  source {
      content = data.template_file.bot_json.rendered
      filename = "QnABot/Bot.json"
  }
  source {
      content = data.template_file.slot_types.rendered
      filename = "QnABot/BotLocales/en_US/SlotTypes/QnaSlotType/SlotType.json"
  }


  # content that is not templated
  source {
      content = file("${path.module}/Manifest.json")
      filename = "Manifest.json"   
  }
  source {
      content = file("${path.module}/QnABot/BotLocales/en_US/BotLocale.json")
      filename = "QnABot/BotLocales/en_US/BotLocale.json"
  }
  source {
      content = file("${path.module}/QnABot/BotLocales/en_US/Intents/FallbackIntent/Intent.json")
      filename = "QnABot/BotLocales/en_US/Intents/FallbackIntent/Intent.json"
  }
  source {
      content = file("${path.module}/QnABot/BotLocales/en_US/Intents/QnaIntent/Slots/qnaslot/Slot.json")
      filename = "QnABot/BotLocales/en_US/Intents/QnaIntent/Slots/qnaslot/Slot.json"
  }
  source {
      content = file("${path.module}/QnABot/BotLocales/en_US/Intents/QnaIntent/Intent.json")
      filename = "QnABot/BotLocales/en_US/Intents/QnaIntent/Intent.json"
  }
  source {
      content = file("${path.module}/QnABot/BotLocales/es_US/BotLocale.json")
      filename = "QnABot/BotLocales/es_US/BotLocale.json"
  }
  source {
      content = file("${path.module}/QnABot/BotLocales/es_US/SlotTypes/QnaSlotType/SlotType.json")
      filename = "QnABot/BotLocales/es_US/SlotTypes/QnaSlotType/SlotType.json"
  }
  source {
      content = file("${path.module}/QnABot/BotLocales/es_US/Intents/FallbackIntent/Intent.json")
      filename = "QnABot/BotLocales/es_US/Intents/FallbackIntent/Intent.json"
  }
  source {
      content = file("${path.module}/QnABot/BotLocales/es_US/Intents/QnaIntent/Slots/qnaslot/Slot.json")
      filename = "QnABot/BotLocales/es_US/Intents/QnaIntent/Slots/qnaslot/Slot.json"
  }
  source {
      content = file("${path.module}/QnABot/BotLocales/es_US/Intents/QnaIntent/Intent.json")
      filename = "QnABot/BotLocales/es_US/Intents/QnaIntent/Intent.json"
  }
  source {
      content = file("${path.module}/QnABot/BotLocales/fr_CA/BotLocale.json")
      filename = "QnABot/BotLocales/fr_CA/BotLocale.json"
  }
  source {
      content = file("${path.module}/QnABot/BotLocales/fr_CA/SlotTypes/QnaSlotType/SlotType.json")
      filename = "QnABot/BotLocales/fr_CA/SlotTypes/QnaSlotType/SlotType.json"
  }
  source {
      content = file("${path.module}/QnABot/BotLocales/fr_CA/Intents/FallbackIntent/Intent.json")
      filename = "QnABot/BotLocales/fr_CA/Intents/FallbackIntent/Intent.json"
  }
  source {
      content = file("${path.module}/QnABot/BotLocales/fr_CA/Intents/QnaIntent/Slots/qnaslot/Slot.json")
      filename = "QnABot/BotLocales/fr_CA/Intents/QnaIntent/Slots/qnaslot/Slot.json"
  }
  source {
      content = file("${path.module}/QnABot/BotLocales/fr_CA/Intents/QnaIntent/Intent.json")
      filename = "QnABot/BotLocales/fr_CA/Intents/QnaIntent/Intent.json"
  }
}

resource "local_file" "question_answer_pairs" {
    content     = jsonencode(local.intents)
    filename = "${path.module}/artifacts/pairs.json"
}

output "archive_path" {
  value = data.archive_file.bot.output_path
}