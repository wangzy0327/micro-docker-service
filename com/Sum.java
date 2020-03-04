package com;
import java.io.*;
import java.util.ArrayList;

public class Sum{
    public static void main(String[] args){
        int sum = 0;
        if(args.length < 2){
            System.out.println("您没有传递任何参数！");
            if(args.length > 0)
                recordResult(args[args.length - 1]);
            return ;
        }
        if(args.length > 2){
            System.out.println("您传递了多余的参数！");
            recordResult(args[args.length - 1]);
            return ;
        }
        for(int i = 0;i < args[0].length();i++){
            if(!Character.isDigit(args[0].charAt(i))){
                System.out.println("您传递的参数非法！请传入整型！");
                recordResult(args[args.length - 1]);
                return ;
            }
        }
        int n = Integer.valueOf(args[0]);
        for(int i = 1;i <= n;i++){
            sum += i;
        }
        System.out.println("sum is "+sum);
        dismissRecord();
        //recordResult(args[args.length - 1]);
        return ;
    }
    private static void dismissRecord(){
        File file = new File("publish/publish.txt");
        String str = null;
        if(file.exists() && file.length() != 0){
            try {
                ArrayList<String> list = new ArrayList<>();
                //BufferedReader是可以按行读取文件
                FileInputStream inputStream = new FileInputStream(file);
                BufferedReader bufferedReader = new BufferedReader(new InputStreamReader(inputStream));

                while((str = bufferedReader.readLine()) != null)
                {
                    if(!str.trim().equals("")) {
                        //System.out.println(str);
                        list.add(str);
                        /**
                         * 可以在添加到list后、根据需求操作文件中第一条数据
                         */
                        //System.out.println("提交");
                    }
                }
                bufferedReader.close();
                inputStream.close();
                recordResult(list.get(0));
                //删除第一行、、写回到文本
                list.remove(0);
                FileOutputStream fileOutputStream = new FileOutputStream(file);
                OutputStreamWriter outputStreamWriter = new OutputStreamWriter(fileOutputStream);
                BufferedWriter bw = new BufferedWriter(outputStreamWriter);

                if(list.size() != 0)
                    for (String string : list) {
                        bw.write(string);
                        bw.newLine();
                        //System.out.println(string);
                    }
                else if(list.size() == 0) {
                    bw.write("");
                }
                bw.flush();
                bw.close();
                //file.delete();
            }catch (Exception e) {
                e.printStackTrace();
            }
        }
    } 
    private static void recordResult(String uuid){
        try {
            File txt = new File("subscribe/subscribe.txt");
            if(!txt.getParentFile().exists()){
                txt.getParentFile().mkdirs();
                txt.createNewFile();
            }
            byte bytes[]= new byte[512];
            bytes=(uuid+"\r\n").getBytes();
            int b=bytes.length;   //是字节的长度，不是字符串的长度
            FileOutputStream fos=new FileOutputStream(txt,true);
            //fos.write(bytes,0,b);
            fos.write(bytes);
            fos.close();
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}
