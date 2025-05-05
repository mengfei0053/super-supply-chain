from docx import Document
import os
import re


def get_docx_paths(model_storage_directory):
    contents = os.listdir(model_storage_directory)
    all_entries = os.listdir(model_storage_directory)

    files = [f for f in all_entries if os.path.isfile(os.path.join(model_storage_directory, f))]

    squared = map(lambda x: os.path.join(model_storage_directory, x), contents)

    # print("文件夹内容:", contents)
    return squared



def read_docx(file_path, outname):
    """读取docx文件内容"""
    doc = Document(file_path)

    print(f"=== 文档 '{file_path}' 内容 ===")

    text = ""
    # 读取所有段落
    for i, paragraph in enumerate(doc.paragraphs, 1):
        # print(f"段落 {i}: {paragraph.text}")
        text += paragraph.text

    matches = re.findall(r'(编号\n*\d+\n*)', text)

    if (matches is not None) and len(matches) > 0 and (matches[0] is not None):
        matches2 = re.findall(r'\d+', matches[0])
        with open(outname, 'a', encoding='utf-8') as f:
            f.write(matches2[0] + '\n')
    else:
        with open(outname, 'a', encoding='utf-8') as f:
            f.write(file_path + '错误' + '\n')







def write_bianhao(dir_path, outname):
    docx_paths = get_docx_paths(dir_path)
    for docx_path in docx_paths:
        read_docx(docx_path, outname)


write_bianhao("/Users/menghongfei/projects/super-supply-chain/py-backend/files/全脂检测报告", "全脂检测.txt")
write_bianhao("/Users/menghongfei/projects/super-supply-chain/py-backend/files/脱脂检测报告", "脱脂检测.txt")





