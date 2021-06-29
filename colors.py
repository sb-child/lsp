import colorama

colorLib = colorama
info = colorama.Fore.LIGHTBLUE_EX
warning = colorama.Fore.LIGHTYELLOW_EX
error = colorama.Fore.LIGHTRED_EX + colorama.Style.BRIGHT
reset = colorama.Fore.RESET + colorama.Style.RESET_ALL + colorLib.Back.RESET
important = colorama.Fore.LIGHTCYAN_EX + colorama.Style.BRIGHT
success = colorama.Fore.LIGHTGREEN_EX + colorama.Style.BRIGHT


def f_info(s):
    return info + str(s) + reset


def f_warning(s):
    return warning + str(s) + reset


def f_error(s):
    return error + str(s) + reset


def f_important(s):
    return important + str(s) + reset


def f_success(s):
    return success + str(s) + reset
