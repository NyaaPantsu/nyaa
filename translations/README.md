# Contributing to translation
You can add your own language support or edit it easily to the website.
## Add a language
To add a language you need to copy /translations/en-us.all.json and translate the strings beside the "translation" key. Do not edit "id" which is the id used to display the translation.

You can also, if you have the website installed, create an empty languageCode.all.json (eg. en-us.all.json) and use the following command:

`cd translations && goi18n -flat false en-us.all.json languageCode.all.json` you need to replace languageCode with the actual language code (eg. en-us)

A new file languageCode.untranslated.json will be created with the new translation strings. Translate them and when it's done, run the following command:

`goi18n -flat=false en-us.all.json languageCode.all.json languageCode.untranslated.json` you need to replace languageCode with the actual language code (eg. en-us)


After creating a new translation, create a new translation string inside "en-us.all.json", like the following:
```
    ...
    },
    {
        "id": "language_(languageCode)_name",
        "translation": "(your language name, in English)"
    },
    ...
```
where languageCode is the newly created ISO code (eg. ja-jp, pt-br).


Before pulling, be sure to delete .unstranslated.json file
## Edit a language
To edit a language you can keep tracking of new strings added to en-us.all.json with the use of git and add the new strings manually to your file.

Or you can also, if you have the website installed, use the following command:
`cd translations && goi18n -flat false en-us.all.json languageCode.all.json` you need to replace languageCode with the actual language code (eg. en-us)

A new file languageCode.untranslated.json will be created with the new translation strings. Translate them and when it's done, run the following command:

`goi18n -flat=false en-us.all.json languageCode.all.json languageCode.untranslated.json` you need to replace languageCode with the actual language code (eg. en-us)

Before pulling, be sure to delete .unstranslated.json file
