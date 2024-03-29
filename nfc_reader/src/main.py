"""
このモジュールは、NFCタグから学生番号と名前を読み取り、それらを表示します。
"""

import os
import threading
import time
from dataclasses import dataclass

import nfc
import requests
from dotenv import load_dotenv


@dataclass
class Configuration:
    """
    アプリケーションの設定を格納するデータクラス。
    """

    nfc_system_code: int
    nfc_service_code: int
    nfc_student_num_block_code: int
    nfc_name_block_code: int


@dataclass
class NfcTagInfo:
    """
    NFCタグの情報を格納するデータクラス。
    """

    idm: bytes
    pmm: bytes
    sys: int
    student_num: str
    name: str


def read_nfc_tag(tag: nfc.tag.Tag, config: Configuration) -> NfcTagInfo:
    """
    NFCタグから学生番号と名前を読み取り
    それらと共にタグのIDm、PMm、システムコードを含むNfcTagInfoオブジェクトを返します。

    Parameters:
    tag (nfc.tag.Tag): 読み取りを行うNFCタグ

    Returns:
    NfcTagInfo: タグのIDm、PMm、システムコード、学生番号、名前を含むオブジェクト
    """
    idm, pmm = tag.polling(system_code=config.nfc_system_code)
    tag.idm, tag.pmm, tag.sys = idm, pmm, config.nfc_system_code
    sc = nfc.tag.tt3.ServiceCode(
        config.nfc_service_code >> 6, config.nfc_service_code & 0x3F
    )

    # student_num
    bc = nfc.tag.tt3.BlockCode(config.nfc_student_num_block_code)
    student_num = tag.read_without_encryption([sc], [bc])
    if isinstance(student_num, str):
        student_num = student_num.encode("shift_jis")
    student_num = student_num.decode("shift_jis").strip("\x00").strip("\x001")
    print("student_number : " + str(student_num))

    # name
    bc = nfc.tag.tt3.BlockCode(config.nfc_name_block_code)
    name = tag.read_without_encryption([sc], [bc])
    if isinstance(name, str):
        name = name.encode("shift_jis")
    name = name.decode("shift_jis").strip("\x00").strip("\x001")
    print("name : " + str(name))

    return NfcTagInfo(idm, pmm, config.nfc_system_code, student_num, name)


def send_request_to_api(nfc_tag_info: NfcTagInfo, unix_timestamp: float) -> None:
    """
    NFCタグの情報とUnixタイムスタンプをAPIに送信します。

    Parameters:
    nfc_tag_info (NfcTagInfo): NFCタグの情報を含むオブジェクト
    unix_timestamp (float): Unixタイムスタンプ（1970年1月1日からの経過秒数）
    """
    url = os.environ["API_URL"]
    headers = {"Content-Type": "application/json"}
    body = {
        "student_number": int(nfc_tag_info.student_num),
        "name": nfc_tag_info.name,
        "timestamp": unix_timestamp,
    }

    response = requests.post(url, json=body, headers=headers, timeout=5)
    print(response.status_code)
    print(response.text)


def on_connect(tag: nfc.tag.Tag) -> bool:
    """
    NFCタグが接続されたときに呼び出される関数。
    タグから学生番号と名前を読み取り
    それらと共にタグのIDm、PMm、システムコードを含むNfcTagInfoオブジェクトを返します。

    Parameters:
    tag (nfc.tag.Tag): 接続されたNFCタグ
    """
    # lord environment
    load_dotenv()
    configuration = Configuration(
        nfc_system_code=int(os.environ["NFC_SYSTEM_CODE"], 16),
        nfc_service_code=int(os.environ["NFC_SERVICE_CODE"], 16),
        nfc_student_num_block_code=int(os.environ["NFC_STUDENT_NUM_BLOCK_CODE"]),
        nfc_name_block_code=int(os.environ["NFC_NAME_BLOCK_CODE"]),
    )

    nfc_tag_info = read_nfc_tag(tag, configuration)
    print(nfc_tag_info)
    current_unix_time = time.time()

    thread = threading.Thread(
        target=send_request_to_api, args=(nfc_tag_info, current_unix_time)
    )
    thread.start()

    return True


def on_release(_tag: nfc.tag.Tag) -> None:
    """
    NFC タグがリリースされたときに呼び出される関数。
    タグがリリースされたことをコンソールに表示します。

    Parameters:
    _tag (nfc.tag.Tag): リリースされた NFC タグ
    """
    print("released")


if __name__ == "__main__":

    rdwr_options = {
        "on-connect": on_connect,
        "on-release": on_release,
    }

    # nfc connect
    with nfc.ContactlessFrontend("usb") as clf:
        while True:
            clf.connect(rdwr=rdwr_options)
