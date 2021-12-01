variable "bot_name" {
  type = string
  description = "The name of the bot"
}

variable "bot_description" {
  type        = string
  description = "The description of the bot"
}

variable "intents" {
  type = list(object({
    id        = string
    questions = list(string)
    answer    = string
  }))
}

locals {
  slot_type_values = flatten([for s in var.intents : [for q in s.questions : { "sampleValue" : { "value" : q }, "synonyms" : null }]])
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
    bot_name        = var.bot_name
    bot_description = var.bot_description
  }
}

data "archive_file" "bot" {
  type        = "zip"
  output_path = "${path.module}/archive/bot.zip"

  # content that is templated
  source {
    content  = data.template_file.bot_json.rendered
    filename = "${var.bot_name}/Bot.json"
  }
  source {
    content  = data.template_file.slot_types.rendered
    filename = "${var.bot_name}/BotLocales/en_US/SlotTypes/QnaSlotType/SlotType.json"
  }


  # content that is not templated
  source {
    content  = file("${path.module}/Manifest.json")
    filename = "Manifest.json"
  }
  source {
    content  = file("${path.module}/QnABot/BotLocales/en_US/BotLocale.json")
    filename = "${var.bot_name}/BotLocales/en_US/BotLocale.json"
  }
  source {
    content  = file("${path.module}/QnABot/BotLocales/en_US/Intents/FallbackIntent/Intent.json")
    filename = "${var.bot_name}/BotLocales/en_US/Intents/FallbackIntent/Intent.json"
  }
  source {
    content  = file("${path.module}/QnABot/BotLocales/en_US/Intents/QnaIntent/Slots/qnaslot/Slot.json")
    filename = "${var.bot_name}/BotLocales/en_US/Intents/QnaIntent/Slots/qnaslot/Slot.json"
  }
  source {
    content  = file("${path.module}/QnABot/BotLocales/en_US/Intents/QnaIntent/Intent.json")
    filename = "${var.bot_name}/BotLocales/en_US/Intents/QnaIntent/Intent.json"
  }
  source {
    content  = file("${path.module}/QnABot/BotLocales/es_US/BotLocale.json")
    filename = "${var.bot_name}/BotLocales/es_US/BotLocale.json"
  }
  source {
    content  = file("${path.module}/QnABot/BotLocales/es_US/SlotTypes/QnaSlotType/SlotType.json")
    filename = "${var.bot_name}/BotLocales/es_US/SlotTypes/QnaSlotType/SlotType.json"
  }
  source {
    content  = file("${path.module}/QnABot/BotLocales/es_US/Intents/FallbackIntent/Intent.json")
    filename = "${var.bot_name}/BotLocales/es_US/Intents/FallbackIntent/Intent.json"
  }
  source {
    content  = file("${path.module}/QnABot/BotLocales/es_US/Intents/QnaIntent/Slots/qnaslot/Slot.json")
    filename = "${var.bot_name}/BotLocales/es_US/Intents/QnaIntent/Slots/qnaslot/Slot.json"
  }
  source {
    content  = file("${path.module}/QnABot/BotLocales/es_US/Intents/QnaIntent/Intent.json")
    filename = "${var.bot_name}/BotLocales/es_US/Intents/QnaIntent/Intent.json"
  }
  source {
    content  = file("${path.module}/QnABot/BotLocales/fr_CA/BotLocale.json")
    filename = "${var.bot_name}/BotLocales/fr_CA/BotLocale.json"
  }
  source {
    content  = file("${path.module}/QnABot/BotLocales/fr_CA/SlotTypes/QnaSlotType/SlotType.json")
    filename = "${var.bot_name}/BotLocales/fr_CA/SlotTypes/QnaSlotType/SlotType.json"
  }
  source {
    content  = file("${path.module}/QnABot/BotLocales/fr_CA/Intents/FallbackIntent/Intent.json")
    filename = "${var.bot_name}/BotLocales/fr_CA/Intents/FallbackIntent/Intent.json"
  }
  source {
    content  = file("${path.module}/QnABot/BotLocales/fr_CA/Intents/QnaIntent/Slots/qnaslot/Slot.json")
    filename = "${var.bot_name}/BotLocales/fr_CA/Intents/QnaIntent/Slots/qnaslot/Slot.json"
  }
  source {
    content  = file("${path.module}/QnABot/BotLocales/fr_CA/Intents/QnaIntent/Intent.json")
    filename = "${var.bot_name}/BotLocales/fr_CA/Intents/QnaIntent/Intent.json"
  }
}

# for debugging
# resource "local_file" "question_answer_pairs" {
#     content     = jsonencode(var.intents)
#     filename = "${path.module}/artifacts/pairs.json"
# }

output "archive_path" {
  value = data.archive_file.bot.output_path
}

output "archive_sha" {
  value = data.archive_file.bot.output_sha
}